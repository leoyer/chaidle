package main

import (
	"zrsf.com/hbase_client/types"
	"encoding/xml"
)

func main() {
	//str := "<REQUEST id='hbaseSave' comment='hbase存储'>	<VERSION>1.0</VERSION>	<TABLE>EINVOICE_INFO</TABLE>	<ROWKEY>05000352333321250057</ROWKEY>	<COLUMN>DATA</COLUMN>	<QUALIFIERS>	<PDF>xxxx</PDF><PDF2>22222</PDF2><PDF3>3333</PDF3></QUALIFIERS></REQUEST>"
	request := &types.SaveRequest{}






	quli := types.StringMap{}
	quli["xxx"] = "aa"
	quli["sss"] = "vc"

	request.Qualifiers = quli
	 datas ,err := xml.Marshal(request)
	if err == nil{
		print(string(datas))
	}


}
