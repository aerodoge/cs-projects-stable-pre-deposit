# cs-projects-stable-pre-deposit


## 流程
- Safe批准（approve）一定数量的USDT给ctStable
- Safe存（deposit）一定数量USDT到ctStable



## 地址
- Safe:`0x28d9464a56129A75f5cEe38651C098A55feB3C11`
- Argus:`0x13623ee9047162a2658c7406a5ac2093c5f75541`
- ctStableUSDT:`0x6503de9FE77d256d9d823f2D335Ce83EcE9E153f`
  - topic0:`0xb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da`



## 执行步骤
- Fork区块
- 部署ACL
- 模拟USDT合约的owner，给safe转点USDT
- 启动事件监听程序
- 发送`setDepositLimits`交易

