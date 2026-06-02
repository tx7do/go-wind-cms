'use client';

import React from 'react';

export interface InfoRow {
    label: string;
    value: React.ReactNode;
}

export interface InfoSection {
    title: string;
    rows: InfoRow[];
}

export interface UserInfoSidebarProps {
    sections: InfoSection[];
}

/**
 * 用户信息侧栏 — 包含基础信息和账号信息
 *
 * 用于用户中心页左侧
 */
const UserInfoSidebar: React.FC<UserInfoSidebarProps> = ({sections}) => {
    return (
        <aside>
            <div className="rounded-xl border border-border bg-card p-6">
                {sections.map((section, sIndex) => (
                    <React.Fragment key={sIndex}>
                        {sIndex > 0 && <div className="my-4 border-t border-border"/>}
                        <div>
                            <h3 className="mb-3 text-sm font-semibold text-foreground">{section.title}</h3>
                            <div className="space-y-2">
                                {section.rows.map((row, rIndex) => (
                                    <div key={rIndex} className="flex justify-between text-sm">
                                        <span className="text-muted-foreground/70">{row.label}:</span>
                                        <span className="font-medium text-foreground">{row.value}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </React.Fragment>
                ))}
            </div>
        </aside>
    );
};

export default UserInfoSidebar;
