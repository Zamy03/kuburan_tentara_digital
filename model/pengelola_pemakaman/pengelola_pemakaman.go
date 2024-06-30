package pengelolapemakaman

type PengelolaPemakaman struct {
    PengelolaId  int    `json:"pengelola_id"`
    NamaLengkap  string `json:"nama_lengkap"`
    Jabatan      string `json:"jabatan"`
    NomorTelepon string `json:"nomor_telepon"`
    Email        string `json:"email"`
    Alamat       string `json:"alamat"`
}
