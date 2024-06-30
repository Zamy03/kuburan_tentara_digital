package plotpemakaman

type PlotPemakaman struct {
    PlotId     int     `json:"plot_id"`
    NomorPlot  string  `json:"nomor_plot"`
    StatusPlot string  `json:"status_plot"` // "tersedia", "terisi", atau "cadangan"
    Keterangan string  `json:"keterangan"`
}
