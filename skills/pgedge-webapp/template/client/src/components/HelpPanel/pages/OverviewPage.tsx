import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import overview from './overview.md?raw';

const OverviewPage = (): React.ReactElement => (
    <ReactMarkdown remarkPlugins={[remarkGfm]}>{overview}</ReactMarkdown>
);

export default OverviewPage;
