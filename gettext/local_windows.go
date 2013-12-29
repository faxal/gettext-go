// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

/*
GetSystemDefaultUILanguage to get the original language of the system,
GetUserDefaultUILanguage to get the current user's selection,
EnumUILanguages to see which languages are available.

#include <iostream>
#include <WinNls.h>
#include <Windows.h>

int main() {
	WCHAR_T localeName[LOCALE_NAME_MAX_LENGTH]={0};
	cout<<"Calling GetUserDefaultLocaleName";
	int ret = GetUserDefaultLocaleName(localeName,LOCALE_NAME_MAX_LENGTH);
	if(ret==0)
		cout<<"Cannot retrieve the default locale name."<<endl;
	else
		wcout<<localeName<<endl;
	return 0;
}
*/
