package testing

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func TestInformer(t *testing.T) { // 创建一个用于停止的channel，用于进程退出前通知informer提前退出。因为informer是一个持久运行的Groutine
	stopper := make(chan struct{}, 2)
	defer close(stopper)

	//	get config
	//config, err := rest.InClusterConfig() // 容器内部获取client配置
	//	指定路径获取client配置
	config, err := clientcmd.BuildConfigFromFlags("", "P:\\pphome\\kubernetes\\.kube\\sun_ub99")
	if err != nil {
		panic(err.Error())
	}

	// create the clientset 通过kubernetes.NewForConfig创建clientset对象。informer需要通过clientset与apiserver进行交互
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//	创建informerFactory
	var informerFactory informers.SharedInformerFactory
	// informers.NewSharedInformerFactory实例化sharedInformer对象
	// 第一个参数是ClientSet
	// 第二个参数是多久同步一次
	informerFactory = informers.NewSharedInformerFactory(clientset, 0)
	//	携带操作选项创建informerFactory
	informerFactory = informers.NewSharedInformerFactoryWithOptions(
		clientset,
		0,
		informers.WithNamespace("default"))
	//	Informer方法可以获得特定资源的informer对象 informer ==> 资源的本地管理器
	informer := informerFactory.Core().V1().Pods().Informer()
	//	注册handler
	_, err = informer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		//	创建资源对象时触发的回调方法
		AddFunc: func(obj interface{}) {
			process(obj, "add")
		},
		//	更新资源对象时触发的回调方法
		UpdateFunc: func(obj interface{}, newObj interface{}) {
			process(newObj, "update")
		},
		//	删除资源对象时触发的回调方法
		DeleteFunc: func(obj interface{}) {
			process(obj, "delete")
		},
	})
	if err != nil {
		panic(err.Error())
	}

	//	开始运行informer对象
	informer.Run(stopper)
}
func process(obj interface{}, operat string) {
	cm := obj.(*v1.Pod)
	fmt.Printf("Pod %s: %s/%s\n", operat, cm.Namespace, cm.Name)
}
