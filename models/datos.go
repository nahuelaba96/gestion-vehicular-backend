package models

type Datos struct {
	ID                 int64   `json:"id"` // autogenerado
	Auto               string  `json:"auto"`
	Patente            string  `json:"patente"`
	TipoBitacora       string  `json:"tipo_bitacora"`
	ComponenteRecambio string  `json:"componente_recambio"`
	ComponenteInstalado string `json:"componente_instalado"`
	Marca              string  `json:"marca"`
	Fecha              string  `json:"fecha"`
	Vendedor           string  `json:"vendedor"`
	Kilometro          int64   `json:"kilometro"`
	Costo              float64 `json:"costo"`
	Nota               string  `json:"nota"`
	FechaProximo       string  `json:"fecha_proximo"`
	KilometrosProximo  int64   `json:"kilometros_proximo"`
}
