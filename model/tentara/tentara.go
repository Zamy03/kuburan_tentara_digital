package tentara

type Tentara struct {
    TentaraId      int    `json:"tentara_id"`
    NamaLengkap    string `json:"nama_lengkap"`
    Pangkat        string `json:"pangkat"`
    TanggalLahir   string `json:"tanggal_lahir"`
    TanggalWafat   string `json:"tanggal_wafat"`
    NomorIdentitas string `json:"nomor_identitas"`
    FotoTentara    string `json:"foto_tentara"`
}
