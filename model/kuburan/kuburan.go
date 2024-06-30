package kuburan

type Kuburan struct {
    KuburanId     int    `json:"kuburan_id"`
    PlotId        int    `json:"plot_id"`
    TentaraId     int    `json:"tentara_id"`
    NomorKuburan  string `json:"nomor_kuburan"`
    TanggalDikubur string `json:"tanggal_dikubur"`
    StatusKuburan string `json:"status_kuburan"` // "terawat" atau "tidak terawat"
    Keterangan    string `json:"keterangan"`
}
