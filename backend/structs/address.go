package structs

type AddressCreateRequest struct {
	RecipientName string `json:"recipient_name" binding:"required" label:"Nama Penerima"`
	Phone         string `json:"phone" binding:"required" label:"Nomor Telepon"`
	AddressLine1  string `json:"address_line1" binding:"required" label:"Alamat Lengkap"`
	AddressLine2  string `json:"address_line2" binding:"required" label:"Detail Alamat"`
	District      string `json:"district" binding:"required" label:"Kecamatan"`
	DistricId     string `json:"distric_id" binding:"required" label:"Kecamatan"`
	City          string `json:"city" binding:"required" label:"Kota/Kabupaten"`
	CityId        string `json:"city_id" binding:"required" label:"Kota/Kabupaten"`
	Province      string `json:"province" binding:"required" label:"Provinsi"`
	ProvinceId    string `json:"province_id" binding:"required" label:"Provinsi"`
	PostalCode    string `json:"postal_code" binding:"required" label:"Kode Pos"`
	IsPrimary     bool   `json:"is_primary" binding:"required" label:"Alamat utama"`
}

type AddressUpdateRequest struct {
	Id            uint   `json:"id" binding:"required" label:"ID Alamat"`
	RecipientName string `json:"recipient_name" binding:"required" label:"Nama Penerima"`
	Phone         string `json:"phone" binding:"required" label:"Nomor Telepon"`
	AddressLine1  string `json:"address_line1" binding:"required" label:"Alamat Lengkap"`
	AddressLine2  string `json:"address_line2" binding:"required" label:"Detail Alamat"`
	District      string `json:"district" binding:"required" label:"Kecamatan"`
	DistricId     string `json:"distric_id" binding:"required" label:"Kecamatan"`
	City          string `json:"city" binding:"required" label:"Kota/Kabupaten"`
	CityId        string `json:"city_id" binding:"required" label:"Kota/Kabupaten"`
	Province      string `json:"province" binding:"required" label:"Provinsi"`
	ProvinceId    string `json:"province_id" binding:"required" label:"Provinsi"`
	PostalCode    string `json:"postal_code" binding:"required" label:"Kode Pos"`
	IsPrimary     bool   `json:"is_primary" binding:"required" label:"Alamat utama"`
}
