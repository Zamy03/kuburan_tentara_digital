package kunjungan

type Kunjungan struct {
    KunjunganId     int    `json:"kunjungan_id"`
    KuburanId       int    `json:"kuburan_id"`
    TanggalKunjungan string `json:"tanggal_kunjungan"`
    NamaPengunjung  string `json:"nama_pengunjung"`
    Hubungan        string `json:"hubungan"`
    Keterangan      string `json:"keterangan"`
}
