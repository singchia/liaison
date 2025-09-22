package dao

import "github.com/singchia/liaison/pkg/liaison/repo/model"

func (d *dao) CreateProxy(proxy *model.Proxy) error {
	return d.getDB().Create(proxy).Error
}

func (d *dao) GetProxyByID(id uint) (*model.Proxy, error) {
	var proxy model.Proxy
	if err := d.getDB().Where("id = ?", id).First(&proxy).Error; err != nil {
		return nil, err
	}
	return &proxy, nil
}

func (d *dao) CountProxies() (int64, error) {
	var count int64
	if err := d.getDB().Model(&model.Proxy{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *dao) ListProxies(page, pageSize int) ([]*model.Proxy, error) {
	db := d.getDB()
	// page & page_size
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	var proxies []*model.Proxy
	if err := db.Find(&proxies).Error; err != nil {
		return nil, err
	}
	return proxies, nil
}

func (d *dao) UpdateProxy(proxy *model.Proxy) error {
	updates := map[string]interface{}{}
	if proxy.Name != "" {
		updates["name"] = proxy.Name
	}
	if proxy.Status != 0 {
		updates["status"] = proxy.Status
	}
	if proxy.Description != "" {
		updates["description"] = proxy.Description
	}
	return d.getDB().Model(&model.Proxy{}).Where("id = ?", proxy.ID).Updates(updates).Error
}

func (d *dao) DeleteProxy(id uint) error {
	return d.getDB().Delete(&model.Proxy{}, id).Error
}
