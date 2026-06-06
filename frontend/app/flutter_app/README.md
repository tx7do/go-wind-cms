# CMS Flutter App

A CMS Flutter project.

## Flutter 常用命令

- `flutter doctor` 查看flutter的状态，查看环境配置是否有问题
- `flutter doctor -v` 查看flutter状态的详细信息
- `flutter --version` 查看Flutter的版本信息
- `flutter build apk` 打包安卓包
- `flutter build apk --target-platform android-arm,android-arm64` 指定编译目标打包Android的APK包
- `flutter build ios` 打包苹果ipa
- `flutter run` 运行项目 默认`--debug`
- `flutter run --profile` 运行线上测试包
- `flutter run --release` 运行线上包
- `flutter channel` 查看flutter 的所有分支
- `flutter channel` stable 切换到具体的分支
- `flutter upgrade` 升级Flutter SDK
- `flutter upgrade --force` 如果升级flutter出现问题可以尝试 强制更新
- `flutter logs` 当链接到某一个设备的时候，通过此命令可以查看到当前设备的log
- `flutter screenshot` 可以截取项目当前屏幕展示的图到项目里
- `flutter clean` 清除Flutter的构建输出文件，比如：构建缓存和构建输出文件（apk、ipa、exe等）
- `flutter analyze` Dart默认的linter配置有点弱, 有很多有问题代码也不报错或警告.
  通过此命令可以应用dart的最佳代码实践,
  对一些不好的代码风格提出警告或者直接报错, 从而提高代码质量
- `flutter attach` 混合开发常用命令，1，首先起项目，运行起整个工程；2，到命令行，打开 flutter_lib
  目录（Flutter
  module工程）；3，输入命令：flutter attach
- `flutter test` 当前项目的单元测试
- `flutter downgrade` 从flutter当前channel 下降到上一个稳定版本
- `flutter install` 直接下载apk到手机上，很方便使用，不用重复build or run
- `flutter create` 创建一个flutter 新项目，比如：`flutter create my_app`
- `flutter pub get` 下载依赖库
- `flutter pub outdated` 查看当前软件包所依赖的每个package，确定哪些package 的依赖项已过时，并为您提供有关如何更新它们的建议。
- `flutter pub upgrade` 用于检索当前 Package 所依赖的其它 Package 的最新版本。如果 pubspec.lock
  文件已经存在，则忽略其保存的版本并以
  pubspec 文件中指定的最新版本为主。如有必要，将会创建或更新该文件。
- `flutter pub upgrade --major-versions` 升级三方库到主版本号的最新版本
- `dart run build_runner build --delete-conflicting-outputs`
- `flutter pub cache clean` 清除掉三方库的缓存
  生成和预编译build脚本；处理输入环境和资源；根据前面的脚本和输入信息，开始正式执行builder生成代码；缓存信息，用于下一回生成代码的时候增量判断使用

## Pod命令

- `pod install` 该命令会根据`Podfile`
  当中的依赖库信息去拉取pod依赖库。首次执行命令的时候，如果ios文件夹中不存在`Podfile.lock`文件，则会自动创建之。
- `pod install --repo-update`
- `pod install --verbose`
- `pod update` 该命令会检查`Podfile.lock`文件中的版本信息，然后进行三方库的更新。
- `pod repo`
- `pod repo update`
- `pod list` 列出项目中的三方库
- `pod spec` 管理spec信息
- `pod env` 打印出pod的环境信息
- `pod cache` 管理pod的缓存
- `pod outdate` 查看三方库是否具有最新的版本可用
- `pod search <name>` 搜索可用的三方库
- `pod init` 在项目中创建`Podfile`

## spider插件

安装spider:

```bash
dart pub global activate spider
```

或者

```bash
brew tap birjuvachhani/spider
brew install spider
```

生成代码:

```bash
spider build
```

## Intl插件

```bash
flutter pub run intl_utils:generate
```

或者

```bash
# 全局安装插件
flutter pub global activate intl_utils

# 运行全局插件进行代码生成
dart pub global run intl_utils:generate
```

## proto插件

安装和更新`protoc-gen-dart`插件

```bash
flutter pub global activate protoc_plugin
```
