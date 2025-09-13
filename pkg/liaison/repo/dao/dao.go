package dao

import (
	"sync"
	"time"

	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Dao 接口定义
type Dao interface {
	// 事务相关方法
	Begin() Dao
	Commit() error
	Rollback() error

	// Edge 相关方法
	GetEdge(id uint64) (*model.Edge, error)
	GetEdgeByAccessKey(accessKey string) (*model.AccessKey, *model.Edge, error)
	CreateEdge(edge *model.Edge) error
	GetEdgeByDeviceID(deviceID uint) (*model.Edge, error)
	ListEdges(page, pageSize int) ([]*model.Edge, error)
	UpdateEdge(edge *model.Edge) error
	UpdateEdgeOnlineStatus(edgeID uint64, onlineStatus model.EdgeOnlineStatus) error
	UpdateEdgeHeartbeatAt(edgeID uint64, heartbeatAt time.Time) error
	UpdateEdgeDeviceID(edgeID uint64, deviceID uint) error
	DeleteEdge(id uint64) error

	// AccessKey 相关方法
	CreateAccessKey(accessKey *model.AccessKey) error
	GetAccessKeyByID(id uint) (*model.AccessKey, error)

	// Device 相关方法
	CreateDevice(device *model.Device) error
	CreateEthernetInterface(iface *model.EthernetInterface) error
	GetDeviceByID(id uint) (*model.Device, error)
	ListDevices(page, pageSize int) ([]*model.Device, error)
	UpdateDevice(device *model.Device) error
	UpdateDeviceUsage(deviceID uint, cpuUsage, memoryUsage, diskUsage float32) error

	// Application 相关方法
	CreateApplication(application *model.Application) error
	GetApplicationByID(id uint) (*model.Application, error)
	ListApplications(query *ListApplicationsQuery) ([]*model.Application, error)
	UpdateApplication(application *model.Application) error
	DeleteApplication(id uint) error

	// Proxy 相关方法
	CreateProxy(proxy *model.Proxy) error
	GetProxyByID(id uint) (*model.Proxy, error)
	ListProxies(page, pageSize int) ([]*model.Proxy, error)
	UpdateProxy(proxy *model.Proxy) error
	DeleteProxy(id uint) error

	// Task 相关方法
	CreateTask(task *model.Task) error
	UpdateTaskStatus(taskID uint, status model.TaskStatus) error
	UpdateTaskResult(taskID uint, status model.TaskStatus, result []byte) error
	UpdateTaskError(taskID uint, error string) error

	// 资源清理
	Close() error
}

type dao struct {
	db *gorm.DB
	tx *gorm.DB     // 事务对象
	mu sync.RWMutex // 保护事务状态的互斥锁

	// config
	config *config.Configuration
}

func NewDao(config *config.Configuration) (Dao, error) {
	db, err := gorm.Open(sqlite.Open(config.Manager.DB), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.Exec("PRAGMA synchronous = OFF;")
	sqlDB.Exec("PRAGMA journal_mode = DELETE;")
	sqlDB.Exec("PRAGMA cache_size = -2000;") // 2MB cache
	sqlDB.Exec("PRAGMA temp_store = MEMORY;")
	sqlDB.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	sqlDB.Exec("PRAGMA mmap_size = 268435456;") // 256MB memory map size
	sqlDB.SetMaxOpenConns(0)

	return &dao{
		db:     db,
		config: config,
	}, nil
}

// Begin 开始事务 - 返回新的事务 DAO 实例
func (d *dao) Begin() Dao {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 总是返回一个新的 DAO 实例，包含事务
	tx := d.db.Begin()
	return &dao{
		db:     d.db,
		tx:     tx,
		config: d.config,
	}
}

// Commit 提交事务 - 线程安全
func (d *dao) Commit() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.tx == nil {
		return nil
	}
	err := d.tx.Commit().Error
	d.tx = nil // 清除事务对象
	return err
}

// Rollback 回滚事务 - 线程安全
func (d *dao) Rollback() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.tx == nil {
		return nil
	}
	err := d.tx.Rollback().Error
	d.tx = nil // 清除事务对象
	return err
}

// getDB 获取当前使用的数据库连接（事务或主连接）- 线程安全
func (d *dao) getDB() *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.tx != nil {
		return d.tx
	}
	return d.db
}

func (d *dao) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

/*
使用示例：

// 并发事务示例
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

// 在业务层使用事务
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

// 或者使用 defer 来确保回滚
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

// 多个事务并发执行
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
*/
