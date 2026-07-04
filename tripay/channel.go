package tripay

const (
	// ChannelPermataVA Permata Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 60 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelPermataVA = "PERMATAVA"

	// ChannelBniVA BNI Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelBniVA = "BNIVA"

	// ChannelBriVA BRI Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 60 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelBriVA = "BRIVA"

	// ChannelMandiriVA Mandiri Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 60 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelMandiriVA = "MANDIRIVA"

	// ChannelBcaVA BCA Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 5.500
	ChannelBcaVA = "BCAVA"

	// ChannelMuamalatVA Muamalat Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 60 menit
	// Max Expired : 180 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelMuamalatVA = "MUAMALATVA"

	// ChannelCimbVA CIMB Niaga Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelCimbVA = "CIMBVA"

	// ChannelBsiVA BSI Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 60 menit
	// Max Expired : 180 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelBsiVA = "BSIVA"

	// ChannelOcbcVA OCBC NISP Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelOcbcVA = "OCBCVA"

	// ChannelDanamonVA Danamon Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelDanamonVA = "DANAMONVA"

	// ChannelOtherBankVA Other Bank Virtual Account
	// Min Amount : Rp 10.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 4.250
	ChannelOtherBankVA = "OTHERBANKVA"

	// ChannelAlfamart Alfamart
	// Min Amount : Rp 10.000
	// Max Amount : Rp 2.500.000
	// Min Expired : 60 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 3.500
	ChannelAlfamart = "ALFAMART"

	// ChannelIndomaret Indomaret
	// Min Amount : Rp 10.000
	// Max Amount : Rp 2.500.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 3.500
	ChannelIndomaret = "INDOMARET"

	// ChannelAlfamidi Alfamidi
	// Min Amount : Rp 5.000
	// Max Amount : Rp 2.500.000
	// Min Expired : 60 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 3.500
	ChannelAlfamidi = "ALFAMIDI"

	// ChannelOvo OVO
	// Min Amount : Rp 1.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 4.320 menit
	// Tipe : REDIRECT
	// Biaya Admin : 3%
	ChannelOvo = "OVO"

	// ChannelQris QRIS by ShopeePay
	// Min Amount : Rp 1.000
	// Max Amount : Rp 5.000.000
	// Min Expired : 10 menit
	// Max Expired : 60 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 750 + 0,7%
	ChannelQris = "QRIS"

	// ChannelQrisc QRIS (Customizable)
	// Min Amount : Rp 1.000
	// Max Amount : Rp 5.000.000
	// Min Expired : 10 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 750 + 0,7%
	ChannelQrisc = "QRISC"

	// ChannelQris2 QRIS
	// Min Amount : Rp 1.000
	// Max Amount : Rp 5.000.000
	// Min Expired : 10 menit
	// Max Expired : 1.440 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 750 + 0,7%
	ChannelQris2 = "QRIS2"

	// ChannelDana DANA
	// Min Amount : Rp 1.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 60 menit
	// Tipe : REDIRECT
	// Biaya Admin : 3%
	ChannelDana = "DANA"

	// ChannelShopeePay ShopeePay
	// Min Amount : Rp 1.000
	// Max Amount : Rp 10.000.000
	// Min Expired : 15 menit
	// Max Expired : 60 menit
	// Tipe : REDIRECT
	// Biaya Admin : 3%
	ChannelShopeePay = "SHOPEEPAY"

	// ChannelQrisShopeePay QRIS Custom by ShopeePay
	// Min Amount : Rp 1.000
	// Max Amount : Rp 5.000.000
	// Min Expired : 10 menit
	// Max Expired : 60 menit
	// Tipe : DIRECT
	// Biaya Admin : Rp 750 + 0,7%
	ChannelQrisShopeePay = "QRIS_SHOPEEPAY"
)
