# nacos_select_destroy_subscribe

### What version of nacos-sdk-go are you using?
v2.2.7

### What version of nacos-sever are you using?
NACOS v2.2.3

### What version of Go are you using (`go version`)?
go 1.22.3

### What operating system (Linux, Windows, …) and version?
mac os

### What did you do?
订阅服务变更后；主动查询了一次服务实例列表

### What did you expect to see?
我没有主动取消订阅服务变更；应该推送服务变动情况。

### What did you see instead?
在查询一次服务实例列表后；订阅服务变更的方法失效了；我所有的订阅服务的callback都无法触发。

# 诉求
1. 理论上没有明显取消订阅实例变化；应该一直保证subscribe有效。
2. 在最新版本 v2.2.9版本中；这个示例中的normal方法也无法成功的订阅服务和拿到callback；不知道是不是我的demo写的有问题。


# 项目说明：
1. 将注册一个名为ServerName的服务；有一个常驻的实例，一个实例会以5s的频率进行上下线；观察订阅服务实例的变化情况。
2. normal 方法是正常工作的订阅
3. destroy 方法中，使用了 SelectInstances 破坏了 Subscribe；导致订阅服务实例出现问题
4. noDestroyAfterSelect 方法中 SelectInstances 未破坏 Subscribe；能够正常接收到服务实例变动


