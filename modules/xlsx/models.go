package xlsx

// Contract 合同模板号
type Contract struct {
	ID         int64  `orm:"name(id);ai" json:"id,omitempty"`            // 自增 ID
	ContractID string `orm:"name(contractid);len(50)" json:"contractid"` // 合同编号
	Created    int64  `orm:"name(created)" json:"created,omitempty"`     // 创建时间
	Company    string `orm:"name(company);len(20)" json:"company"`       // 单位名称
	Edition    string `orm:"name(edition);len(100)" json:"edition"`      // 版本号

	Project     int64  `orm:"name(project);len(100)" json:"project"`         // 项目
	Corporation string `orm:"name(corporation);len(100)" json:"corporation"` // 法人代表
	Identity    string `orm:"name(identity);len(20)" json:"identity"`        // 身份证号
	Address     string `orm:"name(address);len(50)" json:"address"`          // 办公地址
	Mobile      string `orm:"name(mobile);len(20)" json:"mobile"`            // 联系电话
	SpaceNO     string `orm:"name(spaceno);len(20)" json:"spaceno"`          // 工地编号
}

// // Information 乙方信息
// type Information struct {
// 	ID          int64  `orm:"name(id);ai" json:"id"`                        // 自增 ID
// 	Name        int64  `orm:"name(name)" json:"name"`                       // 乙方劳动者
// 	Created     int64  `orm:"name(created)" json:"created,omitempty"`       // 创建时间
// 	Identity    string `orm:"name(identity);len(20)" json:"identity"`       // 身份证号
// 	HomeAddress string `orm:"name(homeaddress);len(50)" json:"homeaddress"` // 户籍住址
// 	Address     string `orm:"name(address);len(50)" json:"address"`         // 现住址
// 	Mobile      string `orm:"name(mobile);len(20)" json:"mobile"`           // 联系电话

// 	Day        string `orm:"name(day);len(50)" json:"day"`               // (日)
// 	MonthWages string `orm:"name(month);len(50)" json:"month"`           // (月)
// 	YearWages  string `orm:"name(year);len(50)" json:"year"`             // (年)
// 	DayWages   string `orm:"name(daywages);len(50)" json:"daywages"`     // 工资(日)
// 	MonthWages string `orm:"name(monthwages);len(50)" json:"monthwages"` // 工资(月)
// 	YearWages  string `orm:"name(yearwages);len(50)" json:"yearwages"`   // 工资(年)
// 	City       string `orm:"name(city);len(50)" json:"city"`             // 工作地点(市)
// 	County     string `orm:"name(county);len(50)" json:"county"`         // 工作地点(县/区)
// 	Company    string `orm:"name(company);len(20)" json:"compan"`        // 部门名称
// 	Profession string `orm:"name(profession);len(20)" json:"profession"` // 职位
// 	Statistics string `orm:"name(statistics);len(20)" json:"statistics"` // 核算间隔(月)
// 	Payment    string `orm:"name(payment);len(20)" json:"payment"`       // 核算后支付日期(日)
// 	LostWages  string `orm:"name(lostwages);len(50)" json:"lostwages"`   // 误工费
// }
