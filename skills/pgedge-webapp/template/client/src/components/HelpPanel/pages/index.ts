import OverviewPage from './OverviewPage';

export interface HelpPage {
    key: string;
    title: string;
    component: React.ComponentType;
}

export const pages: HelpPage[] = [
    { key: 'overview', title: 'Overview', component: OverviewPage },
];
