package utils

/*
#include "C:/Users/midni/AppData/Local/Android/android-ndk-r27/toolchains/llvm/prebuilt/windows-x86_64/sysroot/usr/include/jni.h"
#include <stdlib.h>

// Function to get Android ID
const char* getAndroidID(JNIEnv* env, jobject context) {
    jclass secureClass = (*env)->FindClass(env, "android/provider/Settings$Secure");
    jmethodID getStringMethod = (*env)->GetStaticMethodID(env, secureClass, "getString", "(Landroid/content/ContentResolver;Ljava/lang/String;)Ljava/lang/String;");
    jclass settingsClass = (*env)->FindClass(env, "android/content/Context");
    jmethodID getContentResolverMethod = (*env)->GetMethodID(env, settingsClass, "getContentResolver", "()Landroid/content/ContentResolver;");
    jobject contentResolver = (*env)->CallObjectMethod(env, context, getContentResolverMethod);
    jstring androidIdStr = (*env)->NewStringUTF(env, "android_id");
    jobject androidId = (*env)->CallStaticObjectMethod(env, secureClass, getStringMethod, contentResolver, androidIdStr);
    const char* androidIdCStr = (*env)->GetStringUTFChars(env, (jstring) androidId, 0);
    return androidIdCStr;
}
*/
import "C"
import (
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/driver/desktop"
)

func GetDeviceIdentifier() (string, error) {
	var androidID string
	var err error
	err = driver.RunNative(func(ctx any) error {
		desktopDriver := ctx.(desktop.Driver)
		javaEnv := desktopDriver.(interface{ JNIEnv() *C.JNIEnv }).JNIEnv()
		javaObj := desktopDriver.(interface{ JavaObject() C.jobject }).JavaObject()

		androidID = C.GoString(C.getAndroidID(javaEnv, javaObj))
		return nil
	})
	return androidID, err
}
