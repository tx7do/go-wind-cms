import React from 'react';
import {useI18n} from '@/i18n';
import {Dropdown, Button} from 'antd';
import {GlobalOutlined} from '@ant-design/icons';
import styles from './LocaleSwitcher.module.css';

/**
 * 语言切换选择器组件
 *
 * 通常用于 Header 或导航栏中，供用户切换界面语言
 *
 * @example
 * ```tsx
 * import {LocaleSwitcher} from '@/components/layout/LocaleSwitcher';
 *
 * // 在 Header 组件中使用
 * <header>
 *   <Logo />
 *   <nav>...</nav>
 *   <LocaleSwitcher />
 * </header>
 * ```
 */
export const LocaleSwitcher: React.FC = () => {
    const {locale, changeLocale} = useI18n('menu');

    const localeOptions: { key: string; label: React.ReactNode }[] = [
        {key: 'zh-CN', label: '简体中文'},
        {key: 'en-US', label: 'English'},
    ];

    const handleLanguageChange = ({key}: {key: string}) => {
        changeLocale(key);
    };

    return (
        <Dropdown
            menu={{
                items: localeOptions,
                onClick: handleLanguageChange,
                selectedKeys: [locale]
            }}
            trigger={['click']}
        >
            <Button 
                type="text" 
                icon={<GlobalOutlined />}
                className={styles.langSwitcher}
                aria-label="Language switcher"
            >
                {locale === 'zh-CN' ? '中文' : 'EN'}
            </Button>
        </Dropdown>
    );
};
