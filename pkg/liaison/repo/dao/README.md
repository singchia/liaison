# DAO 层事务管理

本 DAO 层提供了完整的事务支持，让上层业务逻辑可以控制事务的生命周期，并支持并发事务操作。

## 核心特性

1. **接口化设计**: 通过 `Dao` 接口提供统一的事务管理方法
2. **事务支持**: 支持 `Begin()`, `Commit()`, `Rollback()` 操作
3. **并发安全**: 使用互斥锁保护事务状态，支持多 goroutine 并发使用
4. **独立事务**: `Begin()` 返回新的 DAO 实例，每个事务都是独立的
5. **自动回滚**: 提供 `WithTransaction` 辅助函数，自动处理事务回滚
6. **上下文支持**: 支持带上下文的事务操作，可以响应取消信号
7. **并发事务**: 支持多个事务同时执行，每个事务都有独立的 DAO 实例

## 基本用法

### 1. 手动事务管理

```go
func (s *Service) CreateEdgeWithDevice(edge *model.Edge, device *model.Device) error {
	// 开始事务，获得新的事务 DAO
	txDao := s.dao.Begin()
	
	// 创建 Edge
	if err := txDao.CreateEdge(edge); err != nil {
		txDao.Rollback()
		return err
	}
	
	// 创建 Device
	if err := txDao.CreateDevice(device); err != nil {
		txDao.Rollback()
		return err
	}
	
	// 提交事务
	return txDao.Commit()
}
```

### 2. 使用 WithTransaction 辅助函数（推荐）

```go
func (s *Service) CreateEdgeWithDevice(edge *model.Edge, device *model.Device) error {
	return WithTransaction(s.dao, func(txDao Dao) error {
		if err := txDao.CreateEdge(edge); err != nil {
			return err
		}
		return txDao.CreateDevice(device)
	})
}
```

### 3. 带上下文的事务

```go
func (s *Service) CreateEdgeWithDeviceWithContext(ctx context.Context, edge *model.Edge, device *model.Device) error {
	return WithTransactionContext(ctx, s.dao, func(txDao Dao) error {
		if err := txDao.CreateEdge(edge); err != nil {
			return err
		}
		return txDao.CreateDevice(device)
	})
}
```

### 4. 多个事务并发执行

```go
func (s *Service) ProcessMultipleTransactions() error {
	// 事务1：创建 Edge
	txDao1 := s.dao.Begin()
	if err := txDao1.CreateEdge(&model.Edge{Name: "edge1"}); err != nil {
		txDao1.Rollback()
		return err
	}
	
	// 事务2：创建 Device（与事务1并发）
	txDao2 := s.dao.Begin()
	if err := txDao2.CreateDevice(&model.Device{Name: "device1"}); err != nil {
		txDao2.Rollback()
		txDao1.Rollback() // 也要回滚事务1
		return err
	}
	
	// 提交两个事务
	if err := txDao1.Commit(); err != nil {
		txDao2.Rollback()
		return err
	}
	
	return txDao2.Commit()
}
```

## 并发事务

### 1. 基本并发事务

```go
func (s *Service) CreateMultipleEdgesConcurrently(edges []*model.Edge) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(edges))
	
	for _, edge := range edges {
		wg.Add(1)
		go func(e *model.Edge) {
			defer wg.Done()
			
			// 每个 goroutine 都有自己的事务 DAO
			txDao := s.dao.Begin()
			defer func() {
				if r := recover(); r != nil {
					txDao.Rollback()
				}
			}()
			
			if err := txDao.CreateEdge(e); err != nil {
				txDao.Rollback()
				errChan <- err
				return
			}
			
			if err := txDao.Commit(); err != nil {
				errChan <- err
			}
		}(edge)
	}
	
	wg.Wait()
	close(errChan)
	
	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	
	return nil
}
```

### 2. 使用并发事务管理器

```go
func (s *Service) CreateMultipleEdgesWithManager(edges []*model.Edge) error {
	ctm := NewConcurrentTransactionManager(s.dao)
	
	tasks := make([]ConcurrentTask, len(edges))
	for i, edge := range edges {
		tasks[i] = func(txDao Dao) error {
			return txDao.CreateEdge(edge)
		}
	}
	
	return ctm.ExecuteConcurrent(tasks)
}
```

### 3. 限制并发数量

```go
func (s *Service) CreateMultipleEdgesWithLimit(edges []*model.Edge, maxConcurrency int) error {
	ctm := NewConcurrentTransactionManager(s.dao)
	
	tasks := make([]ConcurrentTask, len(edges))
	for i, edge := range edges {
		tasks[i] = func(txDao Dao) error {
			return txDao.CreateEdge(edge)
		}
	}
	
	return ctm.ExecuteConcurrentWithLimit(tasks, maxConcurrency)
}
```

### 4. 收集结果的并发事务

```go
func (s *Service) CreateMultipleEdgesWithResults(edges []*model.Edge) ([]*model.Edge, error) {
	ctm := NewConcurrentTransactionManager(s.dao)
	
	tasks := make([]func(Dao) (interface{}, error), len(edges))
	for i, edge := range edges {
		tasks[i] = func(txDao Dao) (interface{}, error) {
			if err := txDao.CreateEdge(edge); err != nil {
				return nil, err
			}
			return edge, nil
		}
	}
	
	results, err := ctm.ExecuteConcurrentWithResult(tasks)
	if err != nil {
		return nil, err
	}
	
	// 类型转换
	edgeResults := make([]*model.Edge, len(results))
	for i, result := range results {
		edgeResults[i] = result.(*model.Edge)
	}
	
	return edgeResults, nil
}
```

## 事务安全

### 1. 自动回滚

`WithTransaction` 函数会自动处理以下情况：
- 函数返回错误时自动回滚
- 发生 panic 时自动回滚
- 上下文取消时自动回滚

### 2. 独立事务实例

每个 `Begin()` 调用都返回一个新的 DAO 实例，包含独立的事务：
- 不同的事务之间完全隔离
- 支持多个事务同时执行
- 每个事务都有自己的生命周期

### 3. 线程安全

使用 `sync.RWMutex` 保护事务状态：
- 写操作（Begin, Commit, Rollback）使用写锁
- 读操作（getDB）使用读锁
- 支持多 goroutine 并发访问

### 4. 资源清理

使用 `defer` 确保在异常情况下也能正确回滚：

```go
func (s *Service) CreateEdgeWithDeviceSafe(edge *model.Edge, device *model.Device) error {
	txDao := s.dao.Begin()
	defer func() {
		if r := recover(); r != nil {
			txDao.Rollback()
		}
	}()
	
	if err := txDao.CreateEdge(edge); err != nil {
		txDao.Rollback()
		return err
	}
	
	if err := txDao.CreateDevice(device); err != nil {
		txDao.Rollback()
		return err
	}
	
	return txDao.Commit()
}
```

## 最佳实践

1. **优先使用 WithTransaction**: 自动处理回滚，代码更简洁
2. **独立事务**: 每个事务使用独立的 DAO 实例，避免状态污染
3. **错误包装**: 在事务函数中包装错误，提供更多上下文信息
4. **上下文传递**: 对于长时间运行的操作，使用带上下文的事务
5. **避免长事务**: 事务应该尽可能短，避免长时间持有数据库连接
6. **合理使用并发**: 根据业务需求选择合适的并发策略
7. **资源管理**: 确保事务正确提交或回滚，避免资源泄漏

## 注意事项

1. 每个事务 DAO 实例只能使用一次，提交或回滚后就不能再次使用
2. 在事务中进行的查询也会在事务中执行
3. 确保在业务逻辑中正确处理事务的提交和回滚
4. 并发事务会增加数据库连接的使用，注意连接池配置
5. 事务隔离级别由数据库配置决定，注意数据一致性要求 