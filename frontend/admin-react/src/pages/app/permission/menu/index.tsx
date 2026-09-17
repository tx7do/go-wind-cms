import { useRef, useState } from 'react';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import TableExportButton from '@/components/common/TableExportButton';
import { ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag, App, Space } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  NodeExpandOutlined,
  NodeCollapseOutlined,
} from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { permissionservicev1_Menu as Menu } from '@/api/generated/admin/service/v1';
import { PaginationQuery } from '@/core';
import { fetchListMenus, useDeleteMenu } from '@/api/hooks/menu';
import { useProTableScrollY } from '@/hooks/useProTableScrollY';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { getRandomColor } from '@/utils/color';
import {
  getMenuTypeMap,
  getStatusMap,
  MENU_TYPE_STATUS,
  STATUS_TYPE_STATUS,
  normalizeAuthority,
} from './constants';
import MenuDrawer from './components/MenuDrawer';

/**
 * 菜单管理页面
 *
 * 树形表格展示菜单层级结构，支持搜索、CRUD 操作
 */
const MenuManagement = () => {
  const { t } = useTranslation('menu');
  const actionRef = useRef<ActionType>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const containerRef = useRef<HTMLDivElement>(null);
  const tableScrollY = useProTableScrollY(containerRef);

  // Drawer 状态管理
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<'create' | 'edit'>('create');
  const [selectedMenu, setSelectedMenu] = useState<Menu | undefined>();

  // 树数据
  const [treeData, setTreeData] = useState<Menu[]>([]);

  const menuTypeMap = getMenuTypeMap(t);
  const statusMap = getStatusMap(t);

  // 删除操作
  const deleteMutation = useDeleteMenu({
    onSuccess: () => {
      message.success(t('deleteSuccess'));
      actionRef.current?.reload();
      queryClient.invalidateQueries({ queryKey: ['listMenus'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('deleteFailed'));
    },
  });

  // 构建菜单树
  const buildMenuTree = (items: Menu[]): Menu[] => {
    const map = new Map<number, Menu & { children?: Menu[] }>();
    const roots: (Menu & { children?: Menu[] })[] = [];

    // 先创建所有节点的映射
    items.forEach((item) => {
      map.set(item.id as number, { ...item, children: [] });
    });

    // 构建树
    map.forEach((node) => {
      if (node.parentId && map.has(node.parentId as number)) {
        const parent = map.get(node.parentId as number)!;
        if (!parent.children) parent.children = [];
        parent.children.push(node);
      } else {
        roots.push(node);
      }
    });

    // 清理叶子节点的空 children，避免 ProTable 显示无效展开箭头
    const cleanEmpty = (nodes: (Menu & { children?: Menu[] })[]) => {
      nodes.forEach((n) => {
        if (n.children && n.children.length === 0) {
          delete n.children;
        } else if (n.children) {
          cleanEmpty(n.children);
        }
      });
    };
    cleanEmpty(roots);

    return roots;
  };

  // 展开全部/折叠全部
  const [expandedRowKeys, setExpandedRowKeys] = useState<readonly React.Key[]>([]);

  const handleExpandAll = (data: Menu[]) => {
    const keys: React.Key[] = [];
    const collectKeys = (items: Menu[]) => {
      items.forEach((item) => {
        keys.push(item.id as number);
        if ((item as any).children?.length) {
          collectKeys((item as any).children);
        }
      });
    };
    collectKeys(data);
    setExpandedRowKeys(keys);
  };

  const handleCollapseAll = () => {
    setExpandedRowKeys([]);
  };

  // 列配置
  const columns: ProColumns<Menu>[] = [
    {
      title: t('menuName'),
      dataIndex: 'meta.title',
      width: 220,
      fixed: 'left',
      render: (_, record) => {
        const meta = (record as any).meta || {};
        // meta 缺失（如直接 API 创建的菜单）时回退展示 name，避免整列显示 "-"
        const titleText = meta.title || (record as any).name || '-';
        return (
          <Space>
            <span>{titleText}</span>
          </Space>
        );
      },
    },
    {
      title: t('type'),
      dataIndex: 'type',
      width: 95,
      hideInSearch: true,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(MENU_TYPE_STATUS).map(([key, status]) => [
          key,
          { text: menuTypeMap[key as keyof typeof menuTypeMap]?.text || key, status },
        ]),
      ),
      render: (_, record) => {
        const type = record.type as keyof typeof menuTypeMap;
        const config = menuTypeMap[type] || { text: type, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('authority'),
      dataIndex: 'meta.authority',
      hideInSearch: true,
      render: (_, record) => {
        const meta = (record as any).meta || {};
        const authorities = normalizeAuthority(meta.authority);
        if (authorities.length === 0) return '-';
        return (
          <Space wrap size={[4, 4]}>
            {authorities.map((auth) => (
              <Tag key={auth} color={getRandomColor(auth)}>
                {auth}
              </Tag>
            ))}
          </Space>
        );
      },
    },
    {
      title: t('path'),
      dataIndex: 'path',
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('component'),
      dataIndex: 'component',
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      width: 95,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(STATUS_TYPE_STATUS).map(([key, status]) => [
          key,
          { text: statusMap[key as keyof typeof statusMap]?.text || key, status },
        ]),
      ),
      render: (_, record) => {
        const status = record.status as keyof typeof statusMap;
        const config = statusMap[status] || { text: status, color: 'default' };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: t('sortOrder'),
      dataIndex: 'meta.order',
      width: 70,
      hideInSearch: true,
      render: (_, record) => {
        const meta = (record as any).meta || {};
        return meta.order ?? '-';
      },
    },
    {
      title: t('action'),
      valueType: 'option',
      width: 100,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setDrawerMode('edit');
            setSelectedMenu(record);
            setDrawerOpen(true);
          }}
        >
          <EditOutlined />
        </a>,
        <Popconfirm
          key="delete"
          title={t('deleteConfirmTitle')}
          description={t('deleteConfirmDesc', {
            name: (record as any).meta?.title || record.path || '',
          })}
          onConfirm={() => record.id && deleteMutation.mutate({ id: record.id })}
          okText={t('common:button.ok')}
          cancelText={t('common:button.cancel')}
        >
          <a style={{ color: 'var(--ant-color-error)' }}><DeleteOutlined /></a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
        <div ref={containerRef} className="page-container-content">
          <ProTable<Menu>
            actionRef={actionRef}
            columns={columns}
            request={async (params, _sorter, _filter) => {
              try {
                const query = new PaginationQuery({
                  paging: {
                    page: 1,
                    pageSize: 1000,
                  },
                  formValues: {
                    'meta.title': params.menuName,
                    status: params.status,
                  },
                });

                const response = await fetchListMenus(query);
                const items = response.items || [];
                const treeData = buildMenuTree(items as Menu[]);
                setTreeData(treeData);

                return {
                  data: treeData,
                  total: items.length,
                  success: true,
                };
              } catch (error: any) {
                message.error(error.message || t('fetchFailed'));
                return {
                  data: [],
                  total: 0,
                  success: false,
                };
              }
            }}
            rowKey="id"
            search={{
              labelWidth: 'auto',
              defaultCollapsed: false,
            }}
            pagination={false}
            toolBarRender={() => [
              <TableExportButton key="export" fetcher={fetchListMenus} columns={columns} filename="menus" />,
              <Button
                key="create"
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => {
                  setDrawerMode('create');
                  setSelectedMenu(undefined);
                  setDrawerOpen(true);
                }}
              >
                {t('create')}
              </Button>,
              <Button
                key="expandAll"
                icon={<NodeExpandOutlined />}
                onClick={() => handleExpandAll(treeData)}
              >
                {t('expandAll')}
              </Button>,
              <Button
                key="collapseAll"
                icon={<NodeCollapseOutlined />}
                onClick={handleCollapseAll}
              >
                {t('collapseAll')}
              </Button>,

            ]}
            options={{
              density: true,
              fullScreen: true,
              setting: true,
              reload: true,
            }}
            size="middle"
            bordered
            cardBordered={false}
            scroll={{
              y: tableScrollY,
              x: 1100,
            }}
            expandable={{
              expandedRowKeys,
              onExpandedRowsChange: (keys) => setExpandedRowKeys(keys),
            }}
          />
        </div>
      </ContentContainer>

      {/* 菜单编辑/创建 Drawer */}
      <MenuDrawer
        open={drawerOpen}
        mode={drawerMode}
        data={selectedMenu}
        onClose={() => {
          setDrawerOpen(false);
          setSelectedMenu(undefined);
        }}
        onSuccess={() => {
          actionRef.current?.reload();
        }}
      />
    </>
  );
};

export default MenuManagement;
