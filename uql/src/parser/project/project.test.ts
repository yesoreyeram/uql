import { uql } from "../index";

describe("parser", () => {
  describe("project", () => {
    describe("extract", () => {
      it("default", async () => {
        expect(await uql(`project "a"=extract('x=([0-9.]+)',0,"a")`, { data: [{ a: "hello x=45.6|wo" }, { a: "hello x=10|wo" }] })).toStrictEqual([{ a: "x=45.6" }, { a: "x=10" }]);
        expect(await uql(`project "a"=extract('x=([0-9.]+)',1,"a")`, { data: [{ a: "hello x=45.6|wo" }, { a: "hello x=10|wo" }] })).toStrictEqual([{ a: "45.6" }, { a: "10" }]);
        expect(await uql(`project "a"=extract('x=([0-9.]+)',1,"a",'number')`, { data: [{ a: "hello x=45.6|wo" }, { a: "hello x=10|wo" }] })).toStrictEqual([{ a: 45.6 }, { a: 10 }]);
        expect(
          await uql(
            `parse-csv --delimiter "#" --columns "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23"
          | extend "date1"=strcat("1",' ',"2")
          | extend "date1"=trim("date1")
          | extend "date2"=extract('(.)(.*)(.)',2,"date1")
          | extend "date"=todatetime("date2")
          | project "date","level"="5","JVMName"="16","size"="18","status"="20","TimeTaken"="23"`,
            {
              data: `[12/05/2022#22:31:48:265] # [da955ab3-c809-4cf6-ad2c-df3a0209def3]# #DEBUG#test#-#-#-#69920###O(RequestType)#Add(API Method)#JVMName#JVM1#size#78787#status#S##TimeTaken#600
        [12/05/2022#22:31:48:265] # [da955ab3-c899-4de6-ad2c-df3a0209def3]# #DEBUG#test#-#-#-#69920###O(RequestType)#Add(API Method)#JVMName#JVM2#size#789#status#S##TimeTaken#90`,
            }
          )
        ).toStrictEqual([
          { date: new Date("12/05/2022 22:31:48:265"), JVMName: "JVM1", TimeTaken: "600", level: "DEBUG", size: "78787", status: "S" },
          { date: new Date("12/05/2022 22:31:48:265"), JVMName: "JVM2", TimeTaken: "90", level: "DEBUG", size: "789", status: "S" },
        ]);
      });
    });
  });
});
