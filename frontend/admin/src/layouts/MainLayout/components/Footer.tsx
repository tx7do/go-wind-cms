import { usePreferencesStore } from '@/core/preferences/store';

export const AppFooter = () => {
  const { copyright } = usePreferencesStore((state) => state.preferences);

  if (!copyright.enable) return null;

  return (
    <div className="text-center text-gray-400 text-sm py-3">
      © {copyright.date} {copyright.companyName}
      {copyright.icp && (
        <a
          href={copyright.icpLink}
          target="_blank"
          rel="noopener noreferrer"
          className="ml-2 hover:text-blue-600 transition"
        >
          {copyright.icp}
        </a>
      )}
    </div>
  );
};

export default AppFooter;
