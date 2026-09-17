import type {ReactNode} from 'react';
import {PageContainer as ProPageContainer} from '@ant-design/pro-components';

interface PageContainerProps {
    title?: string;
    children: ReactNode;
    extra?: ReactNode;
}

const Index = ({title, children, extra}: PageContainerProps) => {
    return (
        <ProPageContainer
            title={title}
            extra={extra}
        >
            {children}
        </ProPageContainer>
    );
};

export default Index;
