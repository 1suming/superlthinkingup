import create from 'zustand';

interface IMyGlobalInfo{
    isSideNavSticky: boolean; // article 的二级菜单是否固定
    sideNavStickyTop: number;
    setIsSideNavSticky: (isSticky: boolean) => void;
    setSideNavStickyTop: (top: number) => void;
}

const MyGlobalInfoStore = create<IMyGlobalInfo>((set) => ({
    isSideNavSticky: false,
    sideNavStickyTop: 0,
    setIsSideNavSticky: (isSticky) => set({ isSideNavSticky: isSticky }),
    setSideNavStickyTop: (top) => set({ sideNavStickyTop: top }),
}));

export default MyGlobalInfoStore;
