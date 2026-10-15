'use client';

import React, {useEffect} from 'react';

import type {contentservicev1_Page} from '@/api/generated/app/service/v1';
import {getPageTitle} from '@/api/hooks/page';
import SectionRenderer from '@/components/page/SectionRenderer';
import BackButton from '@/components/layout/BackButton';
import SectionContainer from '@/components/layout/SectionContainer';
import Image from '@/components/ui/image';
import {useI18nRouter} from '@/i18n/helpers';

export interface CmsPageViewProps {
    page: contentservicev1_Page;
    /** 返回按钮文案 */
    backLabel?: string;
}

/**
 * CMS Page 前台渲染视图：标题 + 封面 + SectionRenderer。
 * redirectUrl 非空时按页面配置整页跳转（对齐 Page.redirect_url 语义）。
 */
export default function CmsPageView({page, backLabel}: CmsPageViewProps) {
    const router = useI18nRouter();
    const title = getPageTitle(page);
    const cover = page.thumbnail || '';

    useEffect(() => {
        if (page.redirectUrl) {
            if (/^https?:\/\//.test(page.redirectUrl)) {
                window.location.href = page.redirectUrl;
            } else {
                router.replace(page.redirectUrl);
            }
        }
    }, [page.redirectUrl, router]);

    return (
        <div className="w-full">
            {backLabel && (
                <SectionContainer noPadding className="!py-6">
                    <BackButton label={backLabel} onClick={() => router.back()}/>
                </SectionContainer>
            )}

            <article className="mx-auto w-full max-w-[900px] px-8 max-md:px-4">
                {cover && (
                    <div className="relative mb-8 h-[300px] overflow-hidden rounded-xl max-md:h-[200px]">
                        <Image src={cover} alt={title} className="h-full w-full object-cover"/>
                    </div>
                )}
                {title && (
                    <h1 className="mb-8 text-3xl font-bold leading-tight text-foreground max-md:text-2xl">
                        {title}
                    </h1>
                )}
                <SectionRenderer sections={page.sections}/>
            </article>
        </div>
    );
}
