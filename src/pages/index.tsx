import Link from '@docusaurus/Link';
import Translate, { translate } from '@docusaurus/Translate';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import type { ReactNode } from 'react';

import styles from './index.module.css';

function HomepageHeader() {
    const { siteConfig } = useDocusaurusContext();

    return (
        <header className={clsx(styles.heroBanner)}>
            <div className={clsx('container', styles.heroInner)}>
                <div className={styles.heroCopy}>
                    <p className={styles.kicker}>
                        <Translate id="homepage.hero.kicker">Go Telegram-Bot Starter</Translate>
                    </p>
                    <Heading as="h1" className={styles.heroTitle}>
                        {siteConfig.title}
                    </Heading>
                    <p className={styles.heroSubtitle}>
                        <Translate id="homepage.hero.subtitle">
                            A Docusaurus documentation site for a Go Telegram bot template with polling, webhook,
                            health probes, inline demos, and pluggable session storage.
                        </Translate>
                    </p>
                    <div className={styles.buttons}>
                        <Link className="button button--primary button--lg" to="/docs/intro">
                            <Translate id="homepage.hero.primaryButton">Read the Docs</Translate>
                        </Link>
                        <Link className="button button--outline button--lg" to="/docs/architecture/overview">
                            <Translate id="homepage.hero.secondaryButton">View Architecture</Translate>
                        </Link>
                    </div>
                </div>
                <div className={styles.heroPanel}>
                    <div className={styles.panelLabel}>
                        <Translate id="homepage.hero.panelLabel">Template Coverage</Translate>
                    </div>
                    <ul className={styles.panelList}>
                        <li>
                            <Translate id="homepage.hero.scope.image">
                                Local development with polling and production delivery with webhooks
                            </Translate>
                        </li>
                        <li>
                            <Translate id="homepage.hero.scope.video">
                                Command handlers, inline queries, callback queries, and Telegram command sync
                            </Translate>
                        </li>
                        <li>
                            <Translate id="homepage.hero.scope.encryption">
                                Redis-backed session state with automatic in-memory fallback for quick starts
                            </Translate>
                        </li>
                        <li>
                            <Translate id="homepage.hero.scope.ops">
                                Health and readiness endpoints, Docker Compose, and graceful shutdown patterns
                            </Translate>
                        </li>
                    </ul>
                </div>
            </div>
        </header>
    );
}

export default function Home(): ReactNode {
    const { siteConfig } = useDocusaurusContext();

    return (
        <Layout
            title={siteConfig.title}
            description={translate({
                id: 'homepage.layout.description',
                message:
                    'Telegram-Bot Template documentation covering architecture, commands, transport modes, sessions, and operations.',
            })}>
            <HomepageHeader />
            <main>
                <HomepageFeatures />
            </main>
        </Layout>
    );
}
