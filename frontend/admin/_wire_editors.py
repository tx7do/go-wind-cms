# -*- coding: utf-8 -*-

# ===== PageEditPage：renderCustomField + 字段循环 → CustomFieldsEditor =====
path = 'src/pages/app/content/page/edit/PageEditPage.tsx'
src = open(path, encoding='utf-8').read()

old = """  const renderCustomField = (field: any) => {
    const name = field.name ?? '';
    const value = formData.customFields?.[name] ?? '';
    switch (field.type) {
      case 'FIELD_TYPE_RICHTEXT':
        return (
          <Input.TextArea
            rows={2}
            value={value}
            size="small"
            onChange={(e) =>
              updateForm({
                customFields: { ...(formData.customFields ?? {}), [name]: e.target.value },
              })
            }
          />
        );
      default:
        return (
          <Input
            size="small"
            value={value}
            onChange={(e) =>
              updateForm({
                customFields: { ...(formData.customFields ?? {}), [name]: e.target.value },
              })
            }
          />
        );
    }
  };
"""
new = ""
assert old in src, "page renderCustomField"
src = src.replace(old, new, 1)

# 找到字段循环块（Form.Item name={['customFields', ...]} 的 map）整体替换
import re
m = re.search(
    r"\{\s*/\*[^*]*动态字段表单[^*]*\*/.*?\}\)\(\)\s*\}\s*</Form\.Item>\s*\)\s*\}\)?\s*\}?",
    src, re.S)
# 上面这个正则太脆，改用行号定位：直接搜包含 name={['customFields', field.name 的块
lines = src.split('\n')
start = end = None
for i, line in enumerate(lines):
    if "name={['customFields', field.name ?? '']}" in line and start is None:
        # 回溯到该字段循环的起始（Form.Item noStyle shouldUpdate 或 fields.map 附近）
        for j in range(i, max(i - 25, 0), -1):
            if 'fieldDefinitions.map' in lines[j]:
                start = j
                break
    if start is not None and end is None and line.strip() == '}))' and i > start:
        end = i
        break
if start is None or end is None:
    raise SystemExit('page loop block not found')
replacement = """<CustomFieldsEditor
  fields={fieldDefinitions}
  value={formData.customFields}
  onChange={(next) => updateForm({ customFields: next })}
/>"""
lines[start:end + 1] = [replacement]
src = '\n'.join(lines)

src = src.replace("import { fetchListFieldDefinitions,",
                  "import CustomFieldsEditor from '@/components/content/CustomFieldsEditor';\nimport { fetchListFieldDefinitions,")
open(path, 'w', encoding='utf-8', newline='').write(src)
print("page edit ok")

# ===== CategoryEditPage：同样替换 =====
path = 'src/pages/app/content/category/edit/CategoryEditPage.tsx'
src = open(path, encoding='utf-8').read()
old = """          {/* 动态字段表单：由绑定的内容模型字段定义驱动 */}
          <Form.Item
            noStyle
            shouldUpdate={(prev, cur) => prev.contentModelId !== cur.contentModelId}
          >
            {() => {
              return fieldDefinitions.map((field) => (
                <Form.Item
                  key={field.id}
                  name={['customFields', field.name ?? '']}
                  label={field.label || field.name}
                >
                  {renderCustomField(field)}
                </Form.Item>
              ));
            }}
          </Form.Item>"""
new = """          {/* 动态字段表单：由绑定的内容模型字段定义驱动 */}
          <Form.Item
            noStyle
            shouldUpdate={(prev, cur) =>
              prev.contentModelId !== cur.contentModelId ||
              prev.customFields !== cur.customFields
            }
          >
            {() => (
              <CustomFieldsEditor
                fields={fieldDefinitions}
                value={form.getFieldValue('customFields')}
                onChange={(next) => form.setFieldValue('customFields', next)}
              />
            )}
          </Form.Item>"""
assert old in src, "category block"
src = src.replace(old, new, 1)

src = src.replace("import { fetchListFieldDefinitions,",
                  "import CustomFieldsEditor from '@/components/content/CustomFieldsEditor';\nimport { fetchListFieldDefinitions,")

# 删除 CategoryEditPage 的 renderCustomField（可能存在的两种形态）
m = re.search(r"  const renderCustomField = \(field: any\) => \{.*?\n  \};\n\n", src, re.S)
if m:
    src = src.replace(m.group(0), '')
open(path, 'w', encoding='utf-8', newline='').write(src)
print("category edit ok")
