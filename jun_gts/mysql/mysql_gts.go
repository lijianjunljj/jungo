package mysql_gts

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TableName 接口，用于获取表名
type TableName interface {
	TableName() string
}

// MysqlGts MySQL 内存表管理工具，类似 Erlang ETS
type MysqlGts struct {
	db           *gorm.DB                    // GORM 数据库实例
	tables       map[string]*TableCache      // 表缓存，key 为表名
	opLogs       map[string][]*OperationLog  // 操作日志，key 为表名
	mu           sync.RWMutex                // 读写锁
	syncInterval time.Duration               // 同步间隔（秒）
	batchSize    int                         // 每次处理的操作日志数量
	stopChan     chan struct{}               // 停止信号
	syncRunning  bool                        // 同步是否正在运行
}

// TableCache 表缓存结构
type TableCache struct {
	data     map[string]interface{} // 数据缓存，key 为主键 ID
	loaded   bool                   // 是否已加载
	mu       sync.RWMutex           // 表级别的读写锁
	primaryKey string               // 主键字段名，默认为 "ID"
}

// OperationLog 操作日志
type OperationLog struct {
	Operation string                 // 操作类型：insert, update, delete
	Data      map[string]interface{} // 操作数据
	ID        string                 // 主键 ID
	Timestamp time.Time              // 操作时间
}

// NewMysqlGts 创建新的 MysqlGts 实例
func NewMysqlGts(db *gorm.DB) *MysqlGts {
	return &MysqlGts{
		db:           db,
		tables:       make(map[string]*TableCache),
		opLogs:       make(map[string][]*OperationLog),
		syncInterval: 5 * time.Second, // 默认 5 秒同步一次
		batchSize:    100,             // 默认每次处理 100 条
		stopChan:     make(chan struct{}),
		syncRunning:  false,
	}
}

// SetSyncInterval 设置同步间隔（秒）
func (m *MysqlGts) SetSyncInterval(seconds int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.syncInterval = time.Duration(seconds) * time.Second
}

// SetBatchSize 设置每次处理的操作日志数量
func (m *MysqlGts) SetBatchSize(size int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batchSize = size
}

// InitTable 初始化一个或多个表
func (m *MysqlGts) InitTable(models ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, model := range models {
		tableName, err := m.getTableName(model)
		if err != nil {
			return fmt.Errorf("获取表名失败: %v", err)
		}

		// 如果表已存在，跳过
		if _, exists := m.tables[tableName]; exists {
			continue
		}

		// 获取主键字段名
		primaryKey := m.getPrimaryKey(model)

		// 初始化表缓存
		m.tables[tableName] = &TableCache{
			data:       make(map[string]interface{}),
			loaded:     false,
			primaryKey: primaryKey,
		}

		// 初始化操作日志
		m.opLogs[tableName] = make([]*OperationLog, 0)
	}

	return nil
}

// Load 加载整表数据到缓存
func (m *MysqlGts) Load(model interface{}) error {
	tableName, err := m.getTableName(model)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 从数据库加载数据
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 创建切片用于存储查询结果
	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	sliceValue := reflect.New(sliceType)
	slice := sliceValue.Interface()

	// 查询所有数据
	if err := m.db.Find(slice).Error; err != nil {
		return fmt.Errorf("加载表数据失败: %v", err)
	}

	// 将数据存入缓存
	tableCache.mu.Lock()
	defer tableCache.mu.Unlock()

	sliceVal := reflect.ValueOf(slice).Elem()
	for i := 0; i < sliceVal.Len(); i++ {
		item := sliceVal.Index(i).Interface()
		id := m.getIDValue(item, tableCache.primaryKey)
		if id != "" {
			tableCache.data[id] = item
		}
	}

	tableCache.loaded = true
	return nil
}

// First 查询第一条记录
func (m *MysqlGts) First(model interface{}, conds ...interface{}) error {
	tableName, err := m.getTableName(model)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 如果缓存未加载，先加载
	tableCache.mu.RLock()
	loaded := tableCache.loaded
	tableCache.mu.RUnlock()

	if !loaded {
		if err := m.Load(model); err != nil {
			return err
		}
	}

	// 从缓存中查找
	tableCache.mu.RLock()
	defer tableCache.mu.RUnlock()

	// 如果有查询条件，需要匹配
	if len(conds) > 0 {
		// 简单的条件匹配（可以根据需要扩展）
		for _, item := range tableCache.data {
			if m.matchConditions(item, conds) {
				return m.copyModel(item, model)
			}
		}
		return gorm.ErrRecordNotFound
	}

	// 没有条件，返回第一条
	for _, item := range tableCache.data {
		return m.copyModel(item, model)
	}

	return gorm.ErrRecordNotFound
}

// Find 查询多条记录
func (m *MysqlGts) Find(dest interface{}, conds ...interface{}) error {
	tableName, err := m.getTableNameFromDest(dest)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 如果缓存未加载，先加载
	tableCache.mu.RLock()
	loaded := tableCache.loaded
	tableCache.mu.RUnlock()

	if !loaded {
		// 需要从 dest 中获取模型类型
		modelType := reflect.TypeOf(dest).Elem().Elem()
		model := reflect.New(modelType).Interface()
		if err := m.Load(model); err != nil {
			return err
		}
	}

	// 从缓存中查找
	tableCache.mu.RLock()
	defer tableCache.mu.RUnlock()

	// 创建结果切片
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type().Elem()

	var results []interface{}
	for _, item := range tableCache.data {
		if len(conds) == 0 || m.matchConditions(item, conds) {
			results = append(results, item)
		}
	}

	// 将结果复制到 dest
	for _, item := range results {
		newItem := reflect.New(destType).Interface()
		if err := m.copyModel(item, newItem); err != nil {
			continue
		}
		destValue.Set(reflect.Append(destValue, reflect.ValueOf(newItem).Elem()))
	}

	return nil
}

// Update 更新记录
func (m *MysqlGts) Update(model interface{}, values map[string]interface{}) error {
	tableName, err := m.getTableName(model)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 如果缓存未加载，先加载
	tableCache.mu.RLock()
	loaded := tableCache.loaded
	tableCache.mu.RUnlock()

	if !loaded {
		if err := m.Load(model); err != nil {
			return err
		}
	}

	id := m.getIDValue(model, tableCache.primaryKey)
	if id == "" {
		return errors.New("无法获取主键 ID")
	}

	// 更新缓存
	tableCache.mu.Lock()
	defer tableCache.mu.Unlock()

	item, exists := tableCache.data[id]
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// 更新数据
	itemValue := reflect.ValueOf(item).Elem()
	for key, value := range values {
		field := itemValue.FieldByName(m.toCamelCase(key))
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}

	tableCache.data[id] = item

	// 记录操作日志
	m.recordOperation(tableName, "update", id, item)

	return nil
}

// Delete 删除记录
func (m *MysqlGts) Delete(model interface{}, conds ...interface{}) error {
	tableName, err := m.getTableName(model)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 如果缓存未加载，先加载
	tableCache.mu.RLock()
	loaded := tableCache.loaded
	tableCache.mu.RUnlock()

	if !loaded {
		if err := m.Load(model); err != nil {
			return err
		}
	}

	tableCache.mu.Lock()
	defer tableCache.mu.Unlock()

	// 查找要删除的记录
	var toDelete []string
	for id, item := range tableCache.data {
		if len(conds) == 0 || m.matchConditions(item, conds) {
			toDelete = append(toDelete, id)
			// 记录操作日志
			m.recordOperation(tableName, "delete", id, item)
		}
	}

	// 从缓存中删除
	for _, id := range toDelete {
		delete(tableCache.data, id)
	}

	if len(toDelete) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Create 创建新记录
func (m *MysqlGts) Create(model interface{}) error {
	tableName, err := m.getTableName(model)
	if err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	m.mu.RLock()
	tableCache, exists := m.tables[tableName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("表 %s 未初始化，请先调用 InitTable", tableName)
	}

	// 如果缓存未加载，先加载
	tableCache.mu.RLock()
	loaded := tableCache.loaded
	tableCache.mu.RUnlock()

	if !loaded {
		if err := m.Load(model); err != nil {
			return err
		}
	}

	id := m.getIDValue(model, tableCache.primaryKey)
	if id == "" {
		return errors.New("无法获取主键 ID")
	}

	// 添加到缓存
	tableCache.mu.Lock()
	defer tableCache.mu.Unlock()

	// 检查是否已存在
	if _, exists := tableCache.data[id]; exists {
		return errors.New("记录已存在")
	}

	tableCache.data[id] = model

	// 记录操作日志
	m.recordOperation(tableName, "insert", id, model)

	return nil
}

// StartSync 启动定时同步
func (m *MysqlGts) StartSync() {
	m.mu.Lock()
	if m.syncRunning {
		m.mu.Unlock()
		return
	}
	m.syncRunning = true
	m.mu.Unlock()

	go m.syncLoop()
}

// StopSync 停止定时同步
func (m *MysqlGts) StopSync() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.syncRunning {
		return
	}
	m.syncRunning = false
	close(m.stopChan)
	m.stopChan = make(chan struct{})
}

// syncLoop 同步循环
func (m *MysqlGts) syncLoop() {
	ticker := time.NewTicker(m.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.syncToDatabase()
		case <-m.stopChan:
			return
		}
	}
}

// syncToDatabase 同步操作日志到数据库
func (m *MysqlGts) syncToDatabase() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for tableName, logs := range m.opLogs {
		if len(logs) == 0 {
			continue
		}

		// 获取要处理的数量
		processCount := m.batchSize
		if len(logs) < processCount {
			processCount = len(logs)
		}

		// 处理操作日志
		toProcess := logs[:processCount]
		remaining := logs[processCount:]

		// 执行数据库操作
		for _, log := range toProcess {
			if err := m.executeOperation(tableName, log); err != nil {
				// 如果操作失败，可以选择记录错误或重试
				fmt.Printf("同步操作失败 [表: %s, 操作: %s, ID: %s]: %v\n", tableName, log.Operation, log.ID, err)
			}
		}

		// 更新操作日志列表
		m.opLogs[tableName] = remaining
	}
}

// executeOperation 执行单个操作
func (m *MysqlGts) executeOperation(tableName string, log *OperationLog) error {
	tableCache, exists := m.tables[tableName]
	if !exists {
		return fmt.Errorf("表 %s 不存在", tableName)
	}

	switch log.Operation {
	case "insert":
		// 从缓存中获取数据
		tableCache.mu.RLock()
		item, exists := tableCache.data[log.ID]
		tableCache.mu.RUnlock()

		if !exists {
			// 如果缓存中没有，说明可能已经被删除，跳过
			return nil
		}

		// 创建记录
		if err := m.db.Create(item).Error; err != nil {
			return err
		}

	case "update":
		// 从缓存中获取数据
		tableCache.mu.RLock()
		item, exists := tableCache.data[log.ID]
		tableCache.mu.RUnlock()

		if !exists {
			// 如果缓存中没有，说明可能已经被删除，跳过
			return nil
		}

		// 更新记录
		if err := m.db.Save(item).Error; err != nil {
			return err
		}

	case "delete":
		// 对于删除操作，我们需要获取模型类型
		// 由于数据已经从缓存中删除，我们需要从表的其他记录中获取模型类型
		tableCache.mu.RLock()
		var sampleModel interface{}
		// 尝试从缓存中获取一个示例来确定模型类型
		for _, sample := range tableCache.data {
			sampleModel = sample
			break
		}
		tableCache.mu.RUnlock()

		// 如果缓存中没有数据，尝试从操作日志的数据中恢复
		// 但操作日志中存储的是 map，我们需要另一种方式
		if sampleModel == nil {
			// 如果缓存完全为空，我们无法确定模型类型
			// 这种情况下，我们可以尝试使用 GORM 的 Table 方法直接删除
			// 但需要知道表名，我们可以通过 tableName 参数获取
			// 使用原生 SQL 删除
			query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, tableCache.primaryKey)
			if err := m.db.Exec(query, log.ID).Error; err != nil {
				return err
			}
			return nil
		}

		// 创建模型实例并设置主键
		modelType := reflect.TypeOf(sampleModel)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		deleteModel := reflect.New(modelType).Interface()
		if err := m.setIDValue(deleteModel, tableCache.primaryKey, log.ID); err != nil {
			return err
		}

		// 删除记录
		if err := m.db.Delete(deleteModel).Error; err != nil {
			return err
		}
	}

	return nil
}

// recordOperation 记录操作日志
func (m *MysqlGts) recordOperation(tableName, operation, id string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	log := &OperationLog{
		Operation: operation,
		ID:        id,
		Data:      m.modelToMap(data),
		Timestamp: time.Now(),
	}

	m.opLogs[tableName] = append(m.opLogs[tableName], log)
}

// 辅助函数

// getTableName 获取表名
func (m *MysqlGts) getTableName(model interface{}) (string, error) {
	if tn, ok := model.(TableName); ok {
		return tn.TableName(), nil
	}

	// 尝试通过反射调用 TableName 方法
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	method := reflect.ValueOf(model).MethodByName("TableName")
	if method.IsValid() {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0].String(), nil
		}
	}

	return "", errors.New("模型未实现 TableName 方法")
}

// getTableNameFromDest 从目标切片获取表名
func (m *MysqlGts) getTableNameFromDest(dest interface{}) (string, error) {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return "", errors.New("dest 必须是指针类型")
	}

	elemType := destType.Elem()
	if elemType.Kind() != reflect.Slice {
		return "", errors.New("dest 必须是切片类型")
	}

	modelType := elemType.Elem()
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	model := reflect.New(modelType).Interface()
	return m.getTableName(model)
}

// getPrimaryKey 获取主键字段名
func (m *MysqlGts) getPrimaryKey(model interface{}) string {
	// 默认主键为 "ID"
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 检查是否有 ID 字段
	if _, found := modelType.FieldByName("ID"); found {
		return "ID"
	}

	// 可以扩展支持其他主键字段
	return "ID"
}

// getIDValue 获取主键值
func (m *MysqlGts) getIDValue(model interface{}, primaryKey string) string {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	field := modelValue.FieldByName(primaryKey)
	if !field.IsValid() {
		return ""
	}

	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", field.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", field.Int())
	}

	return ""
}

// setIDValue 设置主键值
func (m *MysqlGts) setIDValue(model interface{}, primaryKey, id string) error {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	field := modelValue.FieldByName(primaryKey)
	if !field.IsValid() || !field.CanSet() {
		return errors.New("无法设置主键字段")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(id)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		fmt.Sscanf(id, "%d", &val)
		field.SetUint(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		fmt.Sscanf(id, "%d", &val)
		field.SetInt(val)
	default:
		return errors.New("不支持的主键类型")
	}

	return nil
}

// copyModel 复制模型数据
func (m *MysqlGts) copyModel(src, dst interface{}) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	}

	if srcValue.Type() != dstValue.Type() {
		return errors.New("类型不匹配")
	}

	dstValue.Set(srcValue)
	return nil
}

// matchConditions 匹配查询条件
func (m *MysqlGts) matchConditions(item interface{}, conds []interface{}) bool {
	if len(conds) == 0 {
		return true
	}

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	// 处理 GORM 风格的查询条件
	// 例如: First(&user, "123") 或 First(&user, "id = ?", "123")
	for i := 0; i < len(conds); i++ {
		cond := conds[i]

		// 情况1: 直接传入主键值（字符串或数字）
		if idStr, ok := cond.(string); ok {
			itemID := m.getIDValue(item, "ID")
			if itemID == idStr {
				return true
			}
		} else if idInt, ok := cond.(int); ok {
			itemID := m.getIDValue(item, "ID")
			if itemID == fmt.Sprintf("%d", idInt) {
				return true
			}
		} else if idUint, ok := cond.(uint); ok {
			itemID := m.getIDValue(item, "ID")
			if itemID == fmt.Sprintf("%d", idUint) {
				return true
			}
		}

		// 情况2: 查询条件字符串，例如 "id = ?" 或 "mobile = ?"
		if queryStr, ok := cond.(string); ok {
			// 检查是否包含 "=" 和 "?"
			if containsSubstring(queryStr, "=") && containsSubstring(queryStr, "?") {
				if i+1 >= len(conds) {
					return false
				}
				value := conds[i+1]
				i++ // 跳过下一个参数

				// 提取字段名（简单实现：取 "=" 之前的部分并去除空格）
				fieldName := extractFieldName(queryStr)
				if fieldName != "" {
					field := itemValue.FieldByName(m.toCamelCase(fieldName))
					if field.IsValid() {
						fieldValue := field.Interface()
						if !reflect.DeepEqual(fieldValue, value) {
							return false
						}
					} else {
						return false
					}
				}
			}
		}
	}

	return true
}

// containsSubstring 检查字符串是否包含子串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// extractFieldName 从查询条件中提取字段名
func extractFieldName(query string) string {
	// 例如: "id = ?" -> "id", "mobile = ?" -> "mobile"
	parts := make([]string, 0)
	current := ""
	for _, char := range query {
		if char == ' ' || char == '=' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else if char != '?' {
			current += string(char)
		} else {
			break
		}
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// modelToMap 将模型转换为 map（用于操作日志）
func (m *MysqlGts) modelToMap(model interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	modelType := modelValue.Type()
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)
		if fieldValue.CanInterface() {
			result[field.Name] = fieldValue.Interface()
		}
	}

	return result
}

// toCamelCase 将下划线命名转换为驼峰命名（简单实现）
func (m *MysqlGts) toCamelCase(s string) string {
	// 这里实现简单的转换
	// 如果输入已经是驼峰命名（首字母大写），直接返回
	if len(s) == 0 {
		return s
	}
	
	// 检查是否包含下划线
	hasUnderscore := false
	for _, char := range s {
		if char == '_' {
			hasUnderscore = true
			break
		}
	}
	
	if !hasUnderscore {
		// 如果没有下划线，假设已经是驼峰命名，但需要首字母大写
		if s[0] >= 'a' && s[0] <= 'z' {
			return string(s[0]-32) + s[1:]
		}
		return s
	}
	
	// 转换下划线命名到驼峰命名
	result := ""
	nextUpper := false
	for i, char := range s {
		if char == '_' {
			nextUpper = true
			continue
		}
		if i == 0 || nextUpper {
			// 首字母或下划线后的字母大写
			if char >= 'a' && char <= 'z' {
				result += string(char - 32)
			} else {
				result += string(char)
			}
			nextUpper = false
		} else {
			result += string(char)
		}
	}
	
	return result
}

// GetTableStatus 获取表的状态信息（用于调试和监控）
func (m *MysqlGts) GetTableStatus(tableName string) (loaded bool, recordCount int, pendingOps int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	tableCache, exists := m.tables[tableName]
	if !exists {
		return false, 0, 0
	}
	
	tableCache.mu.RLock()
	loaded = tableCache.loaded
	recordCount = len(tableCache.data)
	tableCache.mu.RUnlock()
	
	pendingOps = len(m.opLogs[tableName])
	
	return loaded, recordCount, pendingOps
}

// ClearCache 清空指定表的缓存（谨慎使用）
func (m *MysqlGts) ClearCache(tableName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	tableCache, exists := m.tables[tableName]
	if !exists {
		return fmt.Errorf("表 %s 不存在", tableName)
	}
	
	tableCache.mu.Lock()
	defer tableCache.mu.Unlock()
	
	tableCache.data = make(map[string]interface{})
	tableCache.loaded = false
	
	return nil
}
