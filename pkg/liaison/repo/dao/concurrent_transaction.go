package dao

import (
	"context"
	"fmt"
	"sync"
)

// ConcurrentTransactionManager 并发事务管理器
type ConcurrentTransactionManager struct {
	dao Dao
}

// NewConcurrentTransactionManager 创建并发事务管理器
func NewConcurrentTransactionManager(dao Dao) *ConcurrentTransactionManager {
	return &ConcurrentTransactionManager{
		dao: dao,
	}
}

// ConcurrentTask 并发任务函数类型
type ConcurrentTask TransactionFunc

// ConcurrentTaskWithContext 带上下文的并发任务函数类型
type ConcurrentTaskWithContext func(context.Context, Dao) error

// ExecuteConcurrent 执行并发事务
func (ctm *ConcurrentTransactionManager) ExecuteConcurrent(tasks []ConcurrentTask) error {
	if len(tasks) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	// 启动所有任务
	for _, task := range tasks {
		wg.Add(1)
		go func(t ConcurrentTask) {
			defer wg.Done()

			// 每个 goroutine 获得独立的事务 DAO
			txDao := ctm.dao.Begin()
			defer func() {
				if r := recover(); r != nil {
					txDao.Rollback()
				}
			}()

			// 执行任务
			if err := t(txDao); err != nil {
				txDao.Rollback()
				errChan <- err
				return
			}

			// 提交事务
			if err := txDao.Commit(); err != nil {
				errChan <- err
			}
		}(task)
	}

	// 等待所有任务完成
	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// 返回第一个错误
	if len(errors) > 0 {
		return fmt.Errorf("concurrent transaction failed: %v", errors[0])
	}

	return nil
}

// ExecuteConcurrentWithContext 执行带上下文的并发事务
func (ctm *ConcurrentTransactionManager) ExecuteConcurrentWithContext(ctx context.Context, tasks []ConcurrentTaskWithContext) error {
	if len(tasks) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	// 启动所有任务
	for _, task := range tasks {
		wg.Add(1)
		go func(t ConcurrentTaskWithContext) {
			defer wg.Done()

			// 每个 goroutine 获得独立的事务 DAO
			txDao := ctm.dao.Begin()
			defer func() {
				if r := recover(); r != nil {
					txDao.Rollback()
				}
			}()

			// 执行任务
			if err := t(ctx, txDao); err != nil {
				txDao.Rollback()
				errChan <- err
				return
			}

			// 提交事务
			if err := txDao.Commit(); err != nil {
				errChan <- err
			}
		}(task)
	}

	// 等待所有任务完成或上下文取消
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		close(errChan)
	}

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// 返回第一个错误
	if len(errors) > 0 {
		return fmt.Errorf("concurrent transaction failed: %v", errors[0])
	}

	return nil
}

// ExecuteConcurrentWithLimit 限制并发数量的并发事务执行
func (ctm *ConcurrentTransactionManager) ExecuteConcurrentWithLimit(tasks []ConcurrentTask, maxConcurrency int) error {
	if len(tasks) == 0 {
		return nil
	}

	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	// 创建任务通道
	taskChan := make(chan ConcurrentTask, len(tasks))
	errChan := make(chan error, len(tasks))

	// 发送所有任务到通道
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskChan {
				// 每个任务获得独立的事务 DAO
				txDao := ctm.dao.Begin()
				defer func() {
					if r := recover(); r != nil {
						txDao.Rollback()
					}
				}()

				// 执行任务
				if err := TransactionFunc(task)(txDao); err != nil {
					txDao.Rollback()
					errChan <- err
					continue
				}

				// 提交事务
				if err := txDao.Commit(); err != nil {
					errChan <- err
				}
			}
		}()
	}

	// 等待所有工作协程完成
	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// 返回第一个错误
	if len(errors) > 0 {
		return fmt.Errorf("concurrent transaction failed: %v", errors[0])
	}

	return nil
}

// ExecuteConcurrentWithResult 执行并发事务并收集结果
func (ctm *ConcurrentTransactionManager) ExecuteConcurrentWithResult(tasks []func(Dao) (interface{}, error)) ([]interface{}, error) {
	if len(tasks) == 0 {
		return nil, nil
	}

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		result interface{}
		err    error
		index  int
	}, len(tasks))

	// 启动所有任务
	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t func(Dao) (interface{}, error)) {
			defer wg.Done()

			// 每个 goroutine 获得独立的事务 DAO
			txDao := ctm.dao.Begin()
			defer func() {
				if r := recover(); r != nil {
					txDao.Rollback()
				}
			}()

			// 执行任务
			result, err := t(txDao)
			if err != nil {
				txDao.Rollback()
				resultChan <- struct {
					result interface{}
					err    error
					index  int
				}{nil, err, idx}
				return
			}

			// 提交事务
			if err := txDao.Commit(); err != nil {
				resultChan <- struct {
					result interface{}
					err    error
					index  int
				}{nil, err, idx}
				return
			}

			resultChan <- struct {
				result interface{}
				err    error
				index  int
			}{result, nil, idx}
		}(i, task)
	}

	// 等待所有任务完成
	wg.Wait()
	close(resultChan)

	// 收集结果
	results := make([]interface{}, len(tasks))
	var errors []error

	for res := range resultChan {
		if res.err != nil {
			errors = append(errors, fmt.Errorf("task %d failed: %w", res.index, res.err))
		} else {
			results[res.index] = res.result
		}
	}

	// 如果有错误，返回第一个错误
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return results, nil
}

/*
使用示例：

// 1. 基本并发事务
func (s *Service) CreateMultipleEdgesConcurrently(edges []*model.Edge) error {
	ctm := NewConcurrentTransactionManager(s.dao)

	tasks := make([]ConcurrentTask, len(edges))
	for i, edge := range edges {
		tasks[i] = func(txDao Dao) error {
			return txDao.CreateEdge(edge)
		}
	}

	return ctm.ExecuteConcurrent(tasks)
}

// 2. 带上下文的并发事务
func (s *Service) CreateMultipleEdgesWithContext(ctx context.Context, edges []*model.Edge) error {
	ctm := NewConcurrentTransactionManager(s.dao)

	tasks := make([]ConcurrentTaskWithContext, len(edges))
	for i, edge := range edges {
		tasks[i] = func(ctx context.Context, txDao Dao) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return txDao.CreateEdge(edge)
			}
		}
	}

	return ctm.ExecuteConcurrentWithContext(ctx, tasks)
}

// 3. 限制并发数量
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

// 4. 收集结果的并发事务
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

// 5. 复杂并发场景
func (s *Service) CreateEdgeWithDevicesConcurrently(edgeData []struct {
	edge   *model.Edge
	device *model.Device
}) error {
	ctm := NewConcurrentTransactionManager(s.dao)

	tasks := make([]ConcurrentTask, len(edgeData))
	for i, data := range edgeData {
		tasks[i] = func(txDao Dao) error {
			// 在同一个事务中创建 Edge 和 Device
			if err := txDao.CreateEdge(data.edge); err != nil {
				return err
			}
			return txDao.CreateDevice(data.device)
		}
	}

	return ctm.ExecuteConcurrent(tasks)
}
*/
