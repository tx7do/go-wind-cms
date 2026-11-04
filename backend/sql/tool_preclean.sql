-- 演示数据注入前的内容清场（dev 库 gwc）。
-- 只动内容侧与演示数据要接管的配置表；登录体系/系统配置/其他业务域一律不碰。
BEGIN;

-- 文章侧（依赖顺序：子表在前）
DELETE FROM comment_likes;
DELETE FROM comments;
DELETE FROM content_references;
DELETE FROM interaction_counters;
DELETE FROM post_versions;
DELETE FROM post_translations;
DELETE FROM post_likes;
DELETE FROM post_watches;
DELETE FROM post_tags;
DELETE FROM post_categories;
DELETE FROM posts;

-- 分类 / 页面 / 区块
DELETE FROM category_translations;
DELETE FROM categories;
DELETE FROM page_versions;
DELETE FROM page_translations;
DELETE FROM pages;
DELETE FROM section_translations;
DELETE FROM sections;

-- 标签
DELETE FROM tag_translations;
DELETE FROM tags;

-- 媒体（e2e 上传的临时资产）
DELETE FROM media_variants;
DELETE FROM media_assets;

-- 导航与站点配置（演示数据自带一套完整的）
DELETE FROM navigation_items;
DELETE FROM navigations;
DELETE FROM site_settings;
DELETE FROM sites;

-- 字典条目 i18n（演示数据自带，字典主表由演示文件 TRUNCATE 重灌）
DELETE FROM sys_dict_entry_i18n;

COMMIT;
