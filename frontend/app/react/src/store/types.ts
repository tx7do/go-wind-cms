export interface IRootState {
    access: IAccessState;
    language: ILanguageState;
    theme: IThemeState;
    user: IUserState;
    loading: ILoadingState;
    authentication: IAuthenticationState;
    category: ICategoryState;
    comment: ICommentState;
    fileTransfer: IFileTransferState;
    navigation: INavigationState;
    page: IPageState;
    post: IPostState;
    tag: ITagState;
    userProfile: IUserProfileState;
}

/**
 * 令牌载荷
 */
export interface TokenPayload {
    /**
     * 令牌值
     */
    value: string;
    /**
     * 令牌过期时间
     */
    expiresAt?: number;
}

export interface IAccessState {
    /**
     * 权限码
     */
    permissions: string[]

    /**
     * 登录 accessToken
     */
    accessToken: TokenPayload | null
    /**
     * 登录 accessToken
     */
    refreshToken: TokenPayload | null

    /**
     * 是否已经检查过权限
     */
    isAccessChecked: boolean
    /**
     * 登录是否过期
     */
    loginExpired: boolean
}

export interface ILanguageState {
    locale: string;
}

/**
 * 主题模式
 * - dark: 暗色模式
 * - light: 亮色模式
 * - system: 跟随系统
 */
export type ThemeMode = 'dark' | 'light' | 'system';

export interface IThemeState {
    mode: ThemeMode;
}

/**
 * 用户信息
 */
export interface IUser {
    /**
     * 用户 id
     */
    id: number;

    /**
     * 用户名
     */
    username: string;

    /**
     * 用户昵称
     */
    nickname: string;

    /**
     * 用户实名
     */
    realname: string;

    /**
     * 头像
     */
    avatar: string;

    /**
     * 用户角色
     */
    roles?: string[];
}

export interface IUserState {
    user: IUser | null;
}

export interface IAuthenticationState {
    loading: boolean;
    error: string | null;
}

export interface ICategoryState {
    categories: any[];
    loading: boolean;
}

export interface ICommentState {
    comments: any[];
    loading: boolean;
}

export interface IFileTransferState {
    uploading: boolean;
    downloading: boolean;
    progress: number;
}

export interface INavigationState {
    items: any[];
    loading: boolean;
}

export interface IPageState {
    pages: any[];
    loading: boolean;
}

export interface IPostState {
    posts: any[];
    loading: boolean;
}

export interface ITagState {
    tags: any[];
    loading: boolean;
}

export interface IUserProfileState {
    profile: any | null;
    loading: boolean;
}

export interface ILoadingState {
    isLoading: boolean;
    error: boolean | null;
}

export type RootState = IRootState;
export type AppDispatch = (...args: any[]) => any;
