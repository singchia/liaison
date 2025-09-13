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

func (d *dao) ListProxies(page, pageSize int) ([]*model.Proxy, error) {
	var proxies []*model.Proxy
	if err := d.getDB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&proxies).Error; err != nil {
		return nil, err
	}
	return proxies, nil
}

func (d *dao) UpdateProxy(proxy *model.Proxy) error {
	return d.getDB().Save(proxy).Error
}

func (d *dao) DeleteProxy(id uint) error {
	return d.getDB().Delete(&model.Proxy{}, id).Error
}
