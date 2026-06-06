// bin/clean_and_gen.dart
// 支持 swagger_parser (Dart 原生) 和 openapi-generator-cli (Java 驱动) 两种生成方式
//
// 用法：
//   dart run bin/clean_and_gen.dart                  # 默认使用 swagger_parser
//   dart run bin/clean_and_gen.dart --openapi        # 使用 openapi-generator-cli (需要 Java)

import 'dart:convert';
import 'dart:io';

import 'package:yaml/yaml.dart';

// ============================================================
// 公共配置
// ============================================================

const yamlPath =
    '../../../backend/app/app/service/cmd/server/assets/openapi.yaml';
const tmpJsonPath = './openapi_clean.json';

// openapi-generator-cli 配置
const jarVersion = '7.4.0';
const outputDir = './lib/generated/api';
const generatorType = 'dart-dio';
const configFile = './generator_config.yaml';

// ============================================================
// 入口
// ============================================================

void main(List<String> args) async {
  final useOpenApi = args.contains('--openapi');

  print('========== [START] API 代码生成 ==========');
  print('[INFO] 生成引擎: ${useOpenApi ? "openapi-generator-cli (Java)" : "swagger_parser (Dart)"}');

  // 1. 校验源文件
  final yamlFile = File(yamlPath);
  if (!await yamlFile.exists()) {
    print('[ERROR] 找不到 openapi.yaml 文件！路径: ${yamlFile.absolute.path}');
    return;
  }

  // 2. 清洗：YAML → 结构化解析 → 递归删除 null → 输出 JSON
  final cleaned = await cleanYamlToJson(yamlFile);
  if (cleaned == null) return;

  // 3. 执行生成
  if (useOpenApi) {
    await generateWithOpenApi(yamlPath);
  } else {
    await generateWithSwaggerParser();
  }

  // 4. 清理临时文件
  final tmpFile = File(tmpJsonPath);
  if (await tmpFile.exists()) await tmpFile.delete();
}

// ============================================================
// 公共：清洗 YAML → JSON（结构化深度 null 清理）
// ============================================================

/// 使用 yaml 包做结构化解析，递归删除所有 null 值和 null 列表元素，
/// 输出为 JSON 格式供 swagger_parser 消费。
///
/// 这比正则清洗可靠得多：无论 YAML 中的 null 以何种形式存在
/// （null、~、空值、破损条目），都会在数据层面被彻底清除。
Future<String?> cleanYamlToJson(File yamlFile) async {
  try {
    print('[INFO] 正在以 UTF-8 读取 YAML 文件...');
    final content = await yamlFile.readAsString(encoding: utf8);

    // 1. 用 yaml 包解析为 YamlMap（结构化）
    final YamlMap yamlRoot = loadYaml(content);

    // 2. 转换为普通 Dart Map/List（递归，同时删除 null）
    final dartObj = _yamlToDart(yamlRoot);
    final cleanedObj = _deepRemoveNulls(dartObj);

    // 3. 修复 swagger_parser 已知 bug：
    //    - requestBody 中 application/json: {}（空 schema）→ 删除整个 requestBody
    //    - responses 中 content: {}（空内容）→ 删除 content 键
    //    这些空结构会让 swagger_parser 的 `as Map<String, dynamic>` 崩溃
    final fixedObj = _fixSwaggerParserBugs(cleanedObj);

    // 4. 输出为 JSON
    final jsonStr = const JsonEncoder.withIndent('  ').convert(fixedObj);
    await File(tmpJsonPath).writeAsString(jsonStr, encoding: utf8);

    print('[DONE] YAML → JSON 清洗完成（null 清理 + swagger_parser bug 修复）');
    return jsonStr;
  } catch (e) {
    print('[ERROR] YAML 清洗失败: $e');
    return null;
  }
}

/// 递归将 YamlMap/YamlList 转换为普通 Map<String, dynamic>/List<dynamic>
dynamic _yamlToDart(dynamic node) {
  if (node is YamlMap) {
    final map = <String, dynamic>{};
    for (final entry in node.entries) {
      final key = entry.key.toString();
      map[key] = _yamlToDart(entry.value);
    }
    return map;
  }
  if (node is YamlList) {
    return node.map(_yamlToDart).toList();
  }
  return node;
}

/// 递归删除所有 null 值和 null 列表元素。
/// 保留空 Map（如 `content: {}`）和空 List（如 `scopes: []`），
/// 因为 OpenAPI 规范中空结构是有意义的。
dynamic _deepRemoveNulls(dynamic obj) {
  if (obj is Map) {
    final result = <String, dynamic>{};
    for (final entry in obj.entries) {
      if (entry.value != null) {
        final cleaned = _deepRemoveNulls(entry.value);
        // 保留 null 会被过滤，但保留 cleaned 本身（即使是空 Map/List）
        result[entry.key.toString()] = cleaned;
      }
    }
    return result;
  }
  if (obj is List) {
    final result = <dynamic>[];
    for (final item in obj) {
      if (item != null) {
        result.add(_deepRemoveNulls(item));
      }
    }
    return result;
  }
  return obj;
}

/// 修复 swagger_parser 已知 bug 的兼容性清洗。
///
/// swagger_parser 内部用 `as Map<String, dynamic>` 强制转换 schema，
/// 遇到以下情况会崩溃：
///   1. requestBody.content.XXX: {}（空 schema）→ 删除整个 requestBody
///   2. responses.XXX.content: {}（空响应体）→ 删除 content 键
dynamic _fixSwaggerParserBugs(dynamic obj) {
  if (obj is Map) {
    final result = <String, dynamic>{};
    for (final entry in obj.entries) {
      final key = entry.key.toString();
      final value = _fixSwaggerParserBugs(entry.value);

      // 跳过无 content 或 content 为空的 requestBody
      // 例如: requestBody: { required: true }  (content 被递归清洗删除了)
      // 例如: requestBody: { content: { application/json: {} } }  (空 schema)
      if (key == 'requestBody' && value is Map) {
        final content = value['content'];
        if (content == null ||
            (content is Map && content.isEmpty) ||
            (content is Map && _allContentTypesEmpty(content))) {
          print('[FIX] 移除空 requestBody（无有效 content/schema）');
          continue;
        }
      }

      // 跳过 responses 中空 content
      // 例如: responses: { "200": { description: "OK", content: {} } }
      if (key == 'content' && value is Map && value.isEmpty) {
        // 空 content，不保留
        continue;
      }

      // 跳过 content type 中空 schema
      // 例如: application/json: {}
      if (value is Map && value.isEmpty && _looksLikeContentType(key)) {
        continue;
      }

      result[key] = value;
    }
    return result;
  }
  if (obj is List) {
    return obj.map(_fixSwaggerParserBugs).toList();
  }
  return obj;
}

/// 检查 content map 中所有 content type 是否都是空 schema
bool _allContentTypesEmpty(Map content) {
  if (content.isEmpty) return true;
  for (final entry in content.entries) {
    final value = entry.value;
    if (value is Map && value.isEmpty) continue;
    // 如果有非空的 content type（包含 schema），则不为空
    return false;
  }
  return true;
}

/// 判断 key 是否像 content type（如 application/json, application/octet-stream）
bool _looksLikeContentType(String key) {
  return key.contains('/') || // application/json, multipart/form-data
      key == 'text/plain' ||
      key == 'text/html';
}

// ============================================================
// 方式一：swagger_parser（Dart 原生，无需 Java）
// ============================================================

Future<void> generateWithSwaggerParser() async {
  try {
    print('[INFO] 正在调用 swagger_parser 解析并生成代码...');

    // 使用清洗后的 JSON 文件（swagger_parser 支持 JSON 格式）
    // 注意：同时指定 --output_directory 确保输出位置，不依赖配置文件
    final parserProcess = await Process.start(
      'dart',
      [
        'run',
        'swagger_parser',
        '--schema_path',
        tmpJsonPath,
        '--output_directory',
        'lib/generated/api',
      ],
      runInShell: true,
    );
    await stdout.addStream(parserProcess.stdout);
    await stderr.addStream(parserProcess.stderr);

    if (await parserProcess.exitCode != 0) {
      print('[ERROR] swagger_parser 解析失败，请检查上方错误。');
      return;
    }

    print('[INFO] 正在调用 build_runner 编译强类型实体 (.g.dart)...');

    final buildProcess = await Process.start(
      'dart',
      ['run', 'build_runner', 'build', '--delete-conflicting-outputs'],
      runInShell: true,
    );
    await stdout.addStream(buildProcess.stdout);
    await stderr.addStream(buildProcess.stderr);

    if (await buildProcess.exitCode == 0) {
      print('========== [DONE] swagger_parser 生成完成！ ==========');
    } else {
      print('[ERROR] build_runner 编译失败。');
    }
  } catch (e) {
    print('[ERROR] swagger_parser 流程异常: $e');
  }
}

// ============================================================
// 方式二：openapi-generator-cli（Java 驱动）
// ============================================================

Future<void> generateWithOpenApi(String sourceYamlPath) async {
  // 1. 确保 jar 文件存在
  final jarFile = await _ensureJarFile();
  if (jarFile == null) return;

  // 2. 校验 Java 环境
  print('[INFO] 正在检查本地 Java 环境...');
  final javaCheck =
      await Process.run('java', ['-version'], runInShell: true);
  if (javaCheck.exitCode != 0) {
    print('[ERROR] 未检测到 Java 环境！此工具底层需要 Java 支持。');
    print('[HINT] 请先安装 JRE 或 JDK，并确保终端能执行 "java" 命令。');
    return;
  }

  // 3. 执行生成
  try {
    print('[INFO] 正在调用 openapi-generator-cli 生成（使用自定义 config）...');

    final process = await Process.start('java', [
      '-jar', jarFile.path,
      'generate',
      '-i', sourceYamlPath,
      '-g', generatorType,
      '-o', outputDir,
      '-c', configFile,
      '--skip-validate-spec',
    ], runInShell: true);

    await stdout.addStream(process.stdout);
    await stderr.addStream(process.stderr);

    final exitCode = await process.exitCode;
    if (exitCode == 0) {
      print('========== [DONE] openapi-generator-cli 生成完成！已导入 $outputDir ==========');
    } else {
      print('[ERROR] 生成失败，请检查上方 Java 的报错日志。');
    }
  } catch (e) {
    print('[ERROR] openapi-generator-cli 流程异常: $e');
  }
}

/// 确保 jar 文件已下载到本地，返回 File 或 null（失败时）
Future<File?> _ensureJarFile() async {
  final jarUrl =
      'https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/$jarVersion/openapi-generator-cli-$jarVersion.jar';

  final toolsDir = Directory('./.tools');
  if (!await toolsDir.exists()) await toolsDir.create();
  final jarFile = File('${toolsDir.path}/openapi-generator-cli.jar');

  // 防御：文件小于 5MB 视为损坏
  if (await jarFile.exists() && await jarFile.length() < 5 * 1024 * 1024) {
    print('[WARN] 检测到本地缓存不完整，正在清理...');
    await jarFile.delete();
  }

  if (await jarFile.exists()) return jarFile;

  // 下载
  print('[INFO] 正在下载 openapi-generator-cli ($jarVersion)... 约 27MB');
  try {
    if (Platform.isWindows) {
      await Process.run('powershell', [
        '-Command',
        "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; "
            "Invoke-WebRequest -Uri '$jarUrl' -OutFile '${jarFile.path}'",
      ], runInShell: true);
    } else {
      await Process.run('curl', [
        '-L', '-o', jarFile.path, jarUrl,
      ], runInShell: true);
    }

    if (await jarFile.exists() && await jarFile.length() > 20 * 1024 * 1024) {
      final mb = (await jarFile.length()) / (1024 * 1024);
      print('[DONE] 核心组件（${mb.toStringAsFixed(2)} MB）下载成功！');
      return jarFile;
    } else {
      print('[ERROR] 下载失败，请确认网络能正常连接 Maven 仓库。');
      if (await jarFile.exists()) await jarFile.delete();
      return null;
    }
  } catch (e) {
    print('[ERROR] 下载异常: $e');
    return null;
  }
}
