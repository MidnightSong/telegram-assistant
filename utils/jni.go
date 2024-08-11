package utils

import "C"
import (
	"fyne.io/fyne/v2/driver"
	"unsafe"
)

/*
#cgo LDFLAGS: -llog
#include <jni.h>
#include <stdlib.h>
#include <android/log.h>

const char* getAndroidID(JNIEnv* env, jobject ctx) {
    jclass secureClass = (*env)->FindClass(env, "android/provider/Settings$Secure");
    jmethodID getStringMethod = (*env)->GetStaticMethodID(env, secureClass, "getString", "(Landroid/content/ContentResolver;Ljava/lang/String;)Ljava/lang/String;");

    jclass contextClass = (*env)->FindClass(env, "android/content/Context");
    jmethodID getContentResolverMethod = (*env)->GetMethodID(env, contextClass, "getContentResolver", "()Landroid/content/ContentResolver;");
    jobject contentResolver = (*env)->CallObjectMethod(env, ctx, getContentResolverMethod);

    jstring androidIDKey = (*env)->NewStringUTF(env, "android_id");
    jstring androidID = (jstring)(*env)->CallStaticObjectMethod(env, secureClass, getStringMethod, contentResolver, androidIDKey);

    const char* id = (*env)->GetStringUTFChars(env, androidID, 0);
    return id;
}
*/
import "C"

func GetAndroidID() (string, error) {
	var androidID string
	err := driver.RunNative(func(ctx any) error {
		if androidCtx, ok := ctx.(unsafe.Pointer); ok {
			// 调用 JNI 获取 ANDROID_ID
			androidID = C.GoString(C.getAndroidID((*C.JNIEnv)(androidCtx), C.jobject(androidCtx)))
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return androidID, nil
}
