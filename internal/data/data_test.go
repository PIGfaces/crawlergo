package data_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/PIGfaces/crawlergo/internal/conf"
	"github.com/PIGfaces/crawlergo/internal/data"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	mockRedisConf = &conf.RedisConf{
		Connection: conf.Connection{
			Host:     "localhost",
			Port:     6379,
			Password: "localtest",
		},
		TargetDB:  1,
		TargetKey: "dd5b0968-e23e-11ec-844f-de8f24aeb194",
		ResultDB:  2,
	}

	testdt = data.NewData(mockRedisConf)
	repo   = data.NewEngineRepo(testdt, mockRedisConf)

	er, _ = repo.(*data.EngineRepo)
)

func setTestData(t *testing.T) {
	// 塞一些测试数据
	testTask := map[string]string{
		uuid.New().String(): "ftp://whlmuxw.mobi/hdesky",
		uuid.New().String(): "wais://odqf.cn/gxryqjged",
		uuid.New().String(): "mid://bldkpnj.pk/yknmjkqrvu",
	}
	testStr, _ := json.Marshal(testTask)
	if err := er.SetRead(context.Background(), mockRedisConf.TargetKey, string(testStr)); err != nil {
		t.Fatal("init test data failed: ", err.Error())
	}
}

func TestNewData(t *testing.T) {
	setTestData(t)
	tasks, err := repo.GetTaskValue(context.Background())
	assert.Nil(t, err, "get task value failed")
	t.Log("===> start mock result")
	for _, task := range tasks {
		result := map[string]interface{}{
			// "原ID:子链接的ID" : { "url": "", "html": "" }
			fmt.Sprintf("%s:%s", task.ID, uuid.New().String()): map[string]string{
				"url":  "prospero://twfzlexngk.ba/hiyn",
				"html": testHtml,
			},
			fmt.Sprintf("%s:%s", task.ID, uuid.New().String()): map[string]string{
				"url":  "http://dummyimage.com/234x60",
				"html": testHtml,
			},
			fmt.Sprintf("%s:%s", task.ID, uuid.New().String()): map[string]string{
				"url":  "uid://ycfdyixq.hk/gpreedlf",
				"html": testHtml,
			},
		}
		for key, value := range result {
			// 序列化
			signalReuslt, err := json.Marshal(value)
			assert.Nil(t, err, "marshal sub task result failed")

			// 保存到缓存中
			err = repo.SetResult(context.Background(), key, string(signalReuslt))
			assert.Nil(t, err, "set redis task failed")

			// 反查结果
			searchResult, err := er.GetResult(context.Background(), key)
			assert.Nil(t, err, "get redis task result failed")

			reverse := map[string]string{}
			// 反序列化
			err = json.Unmarshal(searchResult, &reverse)
			assert.Nil(t, err, "unmarshal redis result failed")

			html, ok := reverse["html"]
			assert.Equal(t, ok, true, "cannot get html from unmarshal map")
			assert.Equal(t, html, testHtml, "set result and get result not equal")
		}
	}
}

const testHtml = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
<!-- 启用Chromium高速渲染模式 -->
<meta name="renderer" content="webkit" />
<meta name="force-rendering" content="webkit"/>
<!-- 禁止百度转码 -->
<meta name="applicable-device" content="pc,mobile" />
<meta name="MobileOptimized" content="width" />
<meta name="HandheldFriendly" content="true" />
<meta http-equiv="Cache-Control" content="no-transform" />
<meta http-equiv="Cache-Control" content="no-siteapp" />
<!-- 禁止识别电话号码 -->
<meta name="format-detection" content="telephone=no" />

<link rel="shortcut icon" href="/favicon.ico?v=1.6.74" />
<link href="/templets/new/style/common.css?v=1.6.74" rel="stylesheet" />
<title>go test命令（Go语言测试命令）完全攻略</title>
<meta name="description" content="Go语言拥有一套单元测试和性能测试系统，仅需要添加很少的代码就可以快速测试一段需求代码。 性能测试系统可以给出代码的性能数据，帮助测试者分析性能问题。 提示 单元测试（" />
</head>
<body>
<div id="topbar" class="clearfix">
	<ul id="product-type" class="left">
		<li>
			<a href="/"><span class="iconfont iconfont-home"></span>首页</a>
		</li>
		<li class="active">
			<a href="/sitemap/" rel="nofollow"><span class="iconfont iconfont-book"></span>教程</a>
		</li>
		<li>
			<a href="http://vip.biancheng.net/p/vip/show.php" rel="nofollow" target="_blank"><span class="iconfont iconfont-vip"></span>VIP会员</a>
		</li>
		<li>
			<a href="http://vip.biancheng.net/p/q2a/show.php" rel="nofollow" target="_blank"><span class="iconfont iconfont-q2a"></span>一对一答疑</a>
		</li>
		<li>
			<a href="http://fudao.biancheng.net/" rel="nofollow" target="_blank"><span class="iconfont iconfont-fudao"></span>辅导班</a>
		</li>
		<!-- <li>
			<a href="/view/9560.html" rel="nofollow" target="_blank"><span class="iconfont iconfont-sword"></span>剑指大厂</a>
		</li> -->
	</ul>
</div>
<div id="header" class="clearfix">
	<a id="logo" class="left" href="/">
		<img height="26" src="/templets/new/images/logo.png?v=1.6.74" alt="C语言中文网" />
	</a>
	<ul id="nav-main" class="hover-none left clearfix">
		<li class="wap-yes"><a href="/">首页</a></li>
		<li><a href="/c/">C语言教程</a></li>
		<li><a href="/cplus/">C++教程</a></li>
		<li><a href="/python/">Python教程</a></li>
		<li><a href="/java/">Java教程</a></li>
		<li><a href="/linux_tutorial/">Linux入门</a></li>
		<li><a href="/sitemap/" title="网站地图">更多&gt;&gt;</a></li>
	</ul>
	<span id="sidebar-toggle" class="toggle-btn" toggle-target="#sidebar">目录 <span class="glyphicon"></span></span>

	<a href="http://vip.biancheng.net/?from=topbar" class="user-info glyphicon glyphicon-user hover-none" target="_blank" rel="nofollow" title="用户中心"></a>
</div>
<div id="main" class="clearfix">
	<div id="sidebar" class="toggle-target">
	<div id="contents">
		<dt><span class="glyphicon glyphicon-option-vertical" aria-hidden="true"></span><a href="/golang/">Go语言</a></dt>
		
			<dd>
				<span class="channel-num">1</span>
				<a href='/golang/intro/'>Go语言简介</a>
			</dd>
		
			<dd>
				<span class="channel-num">2</span>
				<a href='/golang/syntax/'>Go语言基本语法</a>
			</dd>
		
			<dd>
				<span class="channel-num">3</span>
				<a href='/golang/container/'>Go语言容器</a>
			</dd>
		
			<dd>
				<span class="channel-num">4</span>
				<a href='/golang/flow_control/'>流程控制</a>
			</dd>
		
			<dd>
				<span class="channel-num">5</span>
				<a href='/golang/func/'>Go语言函数</a>
			</dd>
		
			<dd>
				<span class="channel-num">6</span>
				<a href='/golang/struct/'>Go语言结构体</a>
			</dd>
		
			<dd>
				<span class="channel-num">7</span>
				<a href='/golang/interface/'>Go语言接口</a>
			</dd>
		
			<dd>
				<span class="channel-num">8</span>
				<a href='/golang/package/'>Go语言包（package）</a>
			</dd>
		
			<dd>
				<span class="channel-num">9</span>
				<a href='/golang/concurrent/'>Go语言并发</a>
			</dd>
		
			<dd>
				<span class="channel-num">10</span>
				<a href='/golang/reflect/'>Go语言反射</a>
			</dd>
		
			<dd>
				<span class="channel-num">11</span>
				<a href='/golang/102/'>Go语言文件处理</a>
			</dd>
		<dd class="this"> <span class="channel-num">12</span> <a href="/golang/build/">Go语言编译与工具</a> </dd><dl class="dl-sub"><dd>12.1 <a href="/view/120.html">go build命令</a></dd><dd>12.2 <a href="/view/4440.html">go clean命令</a></dd><dd>12.3 <a href="/view/121.html">go run命令</a></dd><dd>12.4 <a href="/view/4441.html">go fmt命令</a></dd><dd>12.5 <a href="/view/122.html">go install命令</a></dd><dd>12.6 <a href="/view/123.html">go get命令</a></dd><dd>12.7 <a href="/view/4442.html">go generate命令</a></dd><dd>12.8 <a href="/view/124.html">go test命令</a></dd><dd>12.9 <a href="/view/125.html">go pprof命令</a></dd><dd>12.10 <a href="/view/vip_7361.html">与C/C++进行交互</a><span class="glyphicon glyphicon-usd"></span></dd><dd>12.11 <a href="/view/vip_7362.html">Go语言内存管理</a><span class="glyphicon glyphicon-usd"></span></dd><dd>12.12 <a href="/view/vip_7363.html">Go语言垃圾回收</a><span class="glyphicon glyphicon-usd"></span></dd><dd>12.13 <a href="/view/vip_7365.html">Go语言实现RSA和AES加解密</a><span class="glyphicon glyphicon-usd"></span></dd></dl>
	</div>
</div>
	<div id="article-wrap">
		<div id="article">
			<div class="arc-info">
	<span class="position"><span class="glyphicon glyphicon-map-marker"></span> <a href="/">首页</a> &gt; <a href="/golang/">Go语言</a> &gt; <a href="/golang/build/">Go语言编译与工具</a></span>
	<span class="read-num">阅读：147,530</span>
</div>

<div id="ad-position-bottom"></div>
			<h1>go test命令（Go语言测试命令）完全攻略</h1>
			<div class="pre-next-page clearfix">
                    <span class="pre left"><span class="icon">&lt;</span> <span class="text-brief text-brief-pre">上一页</span><a href="/view/4442.html">go generate命令</a></span>
                    <span class="next right"><a href="/view/125.html">go pprof命令</a><span class="text-brief text-brief-next">下一页</span> <span class="icon">&gt;</span></span>
                </div>
			<div id="ad-arc-top"><p class="pic"></p></div>
			<div id="arc-body">Go语言拥有一套单元测试和性能测试系统，仅需要添加很少的代码就可以快速测试一段需求代码。<br />
<br />
<div>
	go test 命令，会自动读取源码目录下面名为 *_test.go 的文件，生成并运行测试用的可执行文件。输出的信息类似下面所示的样子：</div>
<p class="info-box">
	ok archive/tar 0.011s<br />
	FAIL archive/zip 0.022s<br />
	ok compress/gzip 0.033s<br />
	...</p>
性能测试系统可以给出代码的性能数据，帮助测试者分析性能问题。
<h4>
	提示</h4>
单元测试（unit testing），是指对软件中的最小可测试单元进行检查和验证。对于单元测试中单元的含义，一般要根据实际情况去判定其具体含义，如C语言中单元指一个函数，<a href='/java/' target='_blank'>Java</a> 里单元指一个类，图形化的软件中可以指一个窗口或一个菜单等。总的来说，单元就是人为规定的最小的被测功能模块。<br />
<br />
单元测试是在软件开发过程中要进行的最低级别的测试活动，软件的独立单元将在与程序的其他部分相隔离的情况下进行测试。
<h2>
	单元测试&mdash;&mdash;测试和验证代码的框架</h2>
要开始一个单元测试，需要准备一个 go 源码文件，在命名文件时需要让文件必须以<code style="font-size: 14px;">_test</code>结尾。默认的情况下，<code>go test </code>命令不需要任何的参数，它会自动把你源码包下面所有 test 文件测试完毕，当然你也可以带上参数。<br />
<br />
这里介绍几个常用的参数：
<ul>
	<li>
		-bench regexp 执行相应的 benchmarks，例如 -bench=.；</li>
	<li>
		-cover 开启测试覆盖率；</li>
	<li>
		-run regexp 只运行 regexp 匹配的函数，例如 -run=Array 那么就执行包含有 Array 开头的函数；</li>
	<li>
		-v 显示测试的详细命令。</li>
</ul>
<br />
单元测试源码文件可以由多个测试用例组成，每个测试用例函数需要以<code style="font-size: 14px;">Test</code>为前缀，例如：
<p class="info-box">
	func TestXXX( t *testing.T )</p>
<ul>
	<li>
		测试用例文件不会参与正常源码编译，不会被包含到可执行文件中。</li>
	<li>
		测试用例文件使用<code> go test </code>指令来执行，没有也不需要 main() 作为函数入口。所有在以<code style="font-size: 14px;">_test</code>结尾的源码内以<code style="font-size: 14px;">Test</code>开头的函数会自动被执行。</li>
	<li>
		测试用例可以不传入 *testing.T 参数。</li>
</ul>
<br />
helloworld 的测试代码（具体位置是<code style="font-size: 14px;">./src/chapter11/gotest/helloworld_test.go</code>）：
<blockquote>
	<p>
		本套教程所有源码下载地址：<a href="https://pan.baidu.com/s/1ORFVTOLEYYqDhRzeq0zIiQ" target="_blank">https://pan.baidu.com/s/1ORFVTOLEYYqDhRzeq0zIiQ</a>&nbsp;&nbsp;&nbsp; 提取密码：hfyf</p>
</blockquote>
<pre class="go">
package code11_3

import &quot;testing&quot;

func TestHelloWorld(t *testing.T) {
    t.Log(&quot;hello world&quot;)
}</pre>
代码说明如下：
<ul>
	<li>
		第 5 行，单元测试文件 (*_test.go) 里的测试入口必须以 Test 开始，参数为 *testing.T 的函数。一个单元测试文件可以有多个测试入口。</li>
	<li>
		第 6 行，使用 testing 包的 T 结构提供的 Log() 方法打印字符串。</li>
</ul>
<h4>
	1) 单元测试命令行</h4>
单元测试使用 go test 命令启动，例如：
<pre class="info-box">
$ go test helloworld_test.go
ok          command-line-arguments        0.003s
$ go test -v helloworld_test.go
=== RUN   TestHelloWorld
--- PASS: TestHelloWorld (0.00s)
        helloworld_test.go:8: hello world
PASS
ok          command-line-arguments        0.004s</pre>
代码说明如下：
<ul>
	<li>
		第 1 行，在 go test 后跟 helloworld_test.go 文件，表示测试这个文件里的所有测试用例。</li>
	<li>
		第 2 行，显示测试结果，ok 表示测试通过，command-line-arguments 是测试用例需要用到的一个包名，0.003s 表示测试花费的时间。</li>
	<li>
		第 3 行，显示在附加参数中添加了<code style="font-size: 14px;">-v</code>，可以让测试时显示详细的流程。</li>
	<li>
		第 4 行，表示开始运行名叫 TestHelloWorld 的测试用例。</li>
	<li>
		第 5 行，表示已经运行完 TestHelloWorld 的测试用例，PASS 表示测试成功。</li>
	<li>
		第 6 行打印字符串 hello world。</li>
</ul>
<h4>
	2) 运行指定单元测试用例</h4>
<code>go test </code>指定文件时默认执行文件内的所有测试用例。可以使用<code style="font-size: 14px;">-run</code>参数选择需要的测试用例单独执行，参考下面的代码。<br />
<br />
一个文件包含多个测试用例（具体位置是<code style="font-size: 14px;">./src/chapter11/gotest/select_test.go</code>）
<pre class="go">
package code11_3

import &quot;testing&quot;

func TestA(t *testing.T) {
    t.Log(&quot;A&quot;)
}

func TestAK(t *testing.T) {
    t.Log(&quot;AK&quot;)
}

func TestB(t *testing.T) {
    t.Log(&quot;B&quot;)
}

func TestC(t *testing.T) {
    t.Log(&quot;C&quot;)
}</pre>
这里指定 TestA 进行测试：
<pre class="info-box">
$ go test -v -run TestA select_test.go
=== RUN   TestA
--- PASS: TestA (0.00s)
        select_test.go:6: A
=== RUN   TestAK
--- PASS: TestAK (0.00s)
        select_test.go:10: AK
PASS
ok          command-line-arguments        0.003s</pre>
TestA 和 TestAK 的测试用例都被执行，原因是<code style="font-size: 14px;">-run</code>跟随的测试用例的名称支持正则表达式，使用<code style="font-size: 14px;">-run TestA$</code>即可只执行 TestA 测试用例。
<h4>
	3) 标记单元测试结果</h4>
当需要终止当前测试用例时，可以使用 FailNow，参考下面的代码。<br />
<br />
测试结果标记（具体位置是<code style="font-size: 14px;">./src/chapter11/gotest/fail_test.go</code>）
<pre class="go">
func TestFailNow(t *testing.T) {
    t.FailNow()
}</pre>
还有一种只标记错误不终止测试的方法，代码如下：
<pre class="go">
func TestFail(t *testing.T) {

    fmt.Println(&quot;before fail&quot;)

    t.Fail()

    fmt.Println(&quot;after fail&quot;)
}</pre>
测试结果如下：
<pre class="info-box">
=== RUN   TestFail
before fail
after fail
--- FAIL: TestFail (0.00s)
FAIL
exit status 1
FAIL        command-line-arguments        0.002s</pre>
从日志中看出，第 5 行调用 Fail() 后测试结果标记为失败，但是第 7 行依然被程序执行了。
<h4>
	4) 单元测试日志</h4>
每个测试用例可能并发执行，使用 testing.T 提供的日志输出可以保证日志跟随这个测试上下文一起打印输出。testing.T 提供了几种日志输出方法，详见下表所示。<br />
<br />
<table>
	<caption>
		单元测试框架提供的日志方法</caption>
	<tbody>
		<tr>
			<th>
				方 &nbsp;法</th>
			<th>
				备 &nbsp;注</th>
		</tr>
		<tr>
			<td>
				Log</td>
			<td>
				打印日志，同时结束测试</td>
		</tr>
		<tr>
			<td>
				Logf</td>
			<td>
				格式化打印日志，同时结束测试</td>
		</tr>
		<tr>
			<td>
				Error</td>
			<td>
				打印错误日志，同时结束测试</td>
		</tr>
		<tr>
			<td>
				Errorf</td>
			<td>
				格式化打印错误日志，同时结束测试</td>
		</tr>
		<tr>
			<td>
				Fatal</td>
			<td>
				打印致命日志，同时结束测试</td>
		</tr>
		<tr>
			<td>
				Fatalf</td>
			<td>
				格式化打印致命日志，同时结束测试</td>
		</tr>
	</tbody>
</table>
<br />
开发者可以根据实际需要选择合适的日志。
<h2>
	基准测试&mdash;&mdash;获得代码内存占用和运行效率的性能数据</h2>
基准测试可以测试一段程序的运行性能及耗费 CPU 的程度。Go语言中提供了基准测试框架，使用方法类似于单元测试，使用者无须准备高精度的计时器和各种分析工具，基准测试本身即可以打印出非常标准的测试报告。
<h4>
	1) 基础测试基本使用</h4>
下面通过一个例子来了解基准测试的基本使用方法。<br />
<br />
基准测试（具体位置是<code style="font-size: 14px;">./src/chapter11/gotest/benchmark_test.go</code>）
<pre class="go">
package code11_3

import &quot;testing&quot;

func Benchmark_Add(b *testing.B) {
    var n int
    for i := 0; i &lt; b.N; i++ {
        n++
    }
}</pre>
这段代码使用基准测试框架测试加法性能。第 7 行中的 b.N 由基准测试框架提供。测试代码需要保证函数可重入性及无状态，也就是说，测试代码不使用全局变量等带有记忆性质的<a href='/data_structure/' target='_blank'>数据结构</a>。避免多次运行同一段代码时的环境不一致，不能假设 N 值范围。<br />
<br />
使用如下命令行开启基准测试：
<pre class="info-box">
$ go test -v -bench=. benchmark_test.go
goos: linux
goarch: amd64
Benchmark_Add-4           20000000         0.33 ns/op
PASS
ok          command-line-arguments        0.700s</pre>
代码说明如下：
<ul>
	<li>
		第 1 行的<code style="font-size: 14px;">-bench=.</code>表示运行 benchmark_test.go 文件里的所有基准测试，和单元测试中的<code style="font-size: 14px;">-run</code>类似。</li>
	<li>
		第 4 行中显示基准测试名称，2000000000 表示测试的次数，也就是 testing.B 结构中提供给程序使用的 N。&ldquo;0.33 ns/op&rdquo;表示每一个操作耗费多少时间（纳秒）。</li>
</ul>
<br />
注意：Windows 下使用 go test 命令行时，<code style="font-size: 14px;">-bench=.</code>应写为<code style="font-size: 14px;">-bench=&quot;.&quot;</code>。
<h4>
	2) 基准测试原理</h4>
基准测试框架对一个测试用例的默认测试时间是 1 秒。开始测试时，当以 Benchmark 开头的基准测试用例函数返回时还不到 1 秒，那么 testing.B 中的 N 值将按 1、2、5、10、20、50&hellip;&hellip;递增，同时以递增后的值重新调用基准测试用例函数。
<h4>
	3) 自定义测试时间</h4>
通过<code style="font-size: 14px;">-benchtime</code>参数可以自定义测试时间，例如：
<pre class="info-box">
$ go test -v -bench=. -benchtime=5s benchmark_test.go
goos: linux
goarch: amd64
Benchmark_Add-4           10000000000                 0.33 ns/op
PASS
ok          command-line-arguments        3.380s</pre>
<h4>
	4) 测试内存</h4>
基准测试可以对一段代码可能存在的内存分配进行统计，下面是一段使用字符串格式化的函数，内部会进行一些分配操作。
<pre class="go">
func Benchmark_Alloc(b *testing.B) {

    for i := 0; i &lt; b.N; i++ {
        fmt.Sprintf(&quot;%d&quot;, i)
    }
}</pre>
在命令行中添加<code style="font-size: 14px;">-benchmem</code>参数以显示内存分配情况，参见下面的指令：
<pre class="info-box">
$ go test -v -bench=Alloc -benchmem benchmark_test.go
goos: linux
goarch: amd64
Benchmark_Alloc-4 20000000 109 ns/op 16 B/op 2 allocs/op
PASS
ok          command-line-arguments        2.311s</pre>
代码说明如下：
<ul>
	<li>
		第 1 行的代码中<code style="font-size: 14px;">-bench</code>后添加了 Alloc，指定只测试 Benchmark_Alloc() 函数。</li>
	<li>
		第 4 行代码的&ldquo;16 B/op&rdquo;表示每一次调用需要分配 16 个字节，&ldquo;2 allocs/op&rdquo;表示每一次调用有两次分配。</li>
</ul>
<br />
开发者根据这些信息可以迅速找到可能的分配点，进行优化和调整。
<h4>
	5) 控制计时器</h4>
有些测试需要一定的启动和初始化时间，如果从 Benchmark() 函数开始计时会很大程度上影响测试结果的精准性。testing.B 提供了一系列的方法可以方便地控制计时器，从而让计时器只在需要的区间进行测试。我们通过下面的代码来了解计时器的控制。<br />
<br />
基准测试中的计时器控制（具体位置是<code style="font-size: 14px;">./src/chapter11/gotest/benchmark_test.go</code>）：
<pre class="go">
func Benchmark_Add_TimerControl(b *testing.B) {

    // 重置计时器
    b.ResetTimer()

    // 停止计时器
    b.StopTimer()

    // 开始计时器
    b.StartTimer()

    var n int
    for i := 0; i &lt; b.N; i++ {
        n++
    }
}</pre>
从 Benchmark() 函数开始，Timer 就开始计数。StopTimer() 可以停止这个计数过程，做一些耗时的操作，通过 StartTimer() 重新开始计时。ResetTimer() 可以重置计数器的数据。<br />
<br />
计数器内部不仅包含耗时数据，还包括内存分配的数据。</div>
			<div id="arc-append">
	<p>关注微信公众号「站长严长生」，在手机上阅读所有教程，随时随地都能学习。本公众号由<a class="col-link" href="/view/8092.html" target="_blank" rel="nofollow">C语言中文网站长</a>运营，每日更新，坚持原创，敢说真话，凡事有态度。</p>
	<p style="margin-top:12px; text-align:center;">
		<img width="180" src="/templets/new/images/material/qrcode_weixueyuan_original.png?v=1.6.74" alt="魏雪原二维码"><br>
		<span class="col-green">微信扫描二维码关注公众号</span>
	</p>
</div>
			<div class="pre-next-page clearfix">
                    <span class="pre left"><span class="icon">&lt;</span> <span class="text-brief text-brief-pre">上一页</span><a href="/view/4442.html">go generate命令</a></span>
                    <span class="next right"><a href="/view/125.html">go pprof命令</a><span class="text-brief text-brief-next">下一页</span> <span class="icon">&gt;</span></span>
                </div>
			<div id="ad-arc-bottom"></div>

<!-- <div id="ad-bottom-weixin" class="clearfix">
	<div class="left" style="width: 535px;">
		<p><span class="col-red">编程帮</span>，一个分享编程知识的公众号。跟着<a class="col-link" href="/cpp/about/author/" target="_blank">站长</a>一起学习，每天都有进步。</p>
		<p>通俗易懂，深入浅出，一篇文章只讲一个知识点。</p>
		<p>文章不深奥，不需要钻研，在公交、在地铁、在厕所都可以阅读，随时随地涨姿势。</p>
		<p>文章不涉及代码，不烧脑细胞，人人都可以学习。</p>
		<p>当你决定关注「编程帮」，你已然超越了90%的程序员！</p>
	</div>
	<div class="right" style="width: 150px;">
		<img width="150" src="/templets/new/images/erweima_biancheng.gif?v=1.6.74" alt="编程帮二维码" /><br />
		<span class="col-green">微信扫描二维码关注</span>
	</div>
</div> -->

<div id="nice-arcs" class="box-bottom">
	<h4>优秀文章</h4>
	<ul class="clearfix">
<li><a href="/view/1020.html" title="/etc/rc.d/rc.sysinit配置文件初始化Linux系统">/etc/rc.d/rc.sysinit配置文件初始化Linux系统</a></li>
<li><a href="/view/2016.html" title="C语言二级指针（指向指针的指针）详解">C语言二级指针（指向指针的指针）详解</a></li>
<li><a href="/view/2809.html" title="C# get和set访问器：获取和设置字段（属性）的值">C# get和set访问器：获取和设置字段（属性）的值</a></li>
<li><a href="/view/2812.html" title="C#构造函数（构造方法）">C#构造函数（构造方法）</a></li>
<li><a href="/view/7198.html" title="C++ STL set删除数据：erase()和clear()方法">C++ STL set删除数据：erase()和clear()方法</a></li>
<li><a href="/view/7447.html" title="GDB设置step-mode">GDB设置step-mode</a></li>
<li><a href="/view/7635.html" title="PHP错误类型">PHP错误类型</a></li>
<li><a href="/view/8070.html" title="JS声明变量的3种方式和区别">JS声明变量的3种方式和区别</a></li>
<li><a href="/view/vip_8076.html" title="如何实现C++和C的混合编程？">如何实现C++和C的混合编程？</a></li>
<li><a href="/sql/use.html" title="SQL USE：选择数据库">SQL USE：选择数据库</a></li>
</ul>
</div>
		</div>
		
	</div>
</div>
<script type="text/javascript">
// 当前文章ID
window.arcIdRaw = 'a_' + 124;
window.arcId = "6be87ThatEK0963iaKit3FSTeP8/DJkbWYb76YtH4NrUAxBgafe8wAyI4Q";
window.typeidChain = "14|1";
</script>
<div id="footer" class="clearfix">
	<div class="info left">
	<p>精美而实用的网站，分享优质编程教程，帮助有志青年。千锤百炼，只为大作；精益求精，处处斟酌；这种教程，看一眼就倾心。</p>
	<p>
		<a href="/view/8066.html" target="_blank" rel="nofollow">关于网站</a> <span>|</span>
		<a href="/view/8092.html" target="_blank" rel="nofollow">关于站长</a> <span>|</span>
		<a href="/view/8097.html" target="_blank" rel="nofollow">如何完成一部教程</a> <span>|</span>
		<a href="/view/8093.html" target="_blank" rel="nofollow">联系我们</a> <span>|</span>
		<a href="/sitemap/" target="_blank" rel="nofollow">网站地图</a>
	</p>
	<p>Copyright ©2012-2022 biancheng.net, <a href="http://www.beian.miit.gov.cn/" target="_blank" rel="nofollow" style="color:#666;">陕ICP备15000209号</a></p>
	</div>
	<img class="right" src="/templets/new/images/logo_bottom.gif?v=1.6.74" alt="底部Logo" />
	<span id="return-top"><b>↑</b></span>
</div>

<script type="text/javascript">
window.siteId = 4;
window.cmsTemplets = "/templets/new";
window.cmsTempletsVer = "1.6.74";
</script>

<script src="/templets/new/script/jquery1.12.4.min.js"></script>
<script src="/templets/new/script/common.js?v=1.6.74"></script>
<!-- cnzz统计 -->
<!-- <span style="display:none;"><script src="http://s19.cnzz.com/z_stat.php?id=1274766082&web_id=1274766082" type="text/javascript" defer="defer" async="async"></script></span> -->
<!-- 百度统计 -->
<script>
var _hmt = _hmt || [];
(function() {
var hm = document.createElement("script");
hm.src = "https://hm.baidu.com/hm.js?ef8113430e2dcae44abc0e42dc639917";
var s = document.getElementsByTagName("script")[0]; 
s.parentNode.insertBefore(hm, s);
})();
</script>
</body>
</html>`
