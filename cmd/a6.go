package main

// A6使用的本地资产管理系统上导出的总资产，csv格式
type A6Line struct {
	Company      string `csv:"公司"`
	AssetTag     string `csv:"资产标签"`
	Model        string `csv:"型号"`
	Manufacturer string `csv:"生产厂家"`
}


