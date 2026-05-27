import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';
import clsx from 'clsx';
import type { ReactNode } from 'react';
import styles from './styles.module.css';

type FeatureItem = {
    eyebrowId: string;
    eyebrow: string;
    titleId: string;
    title: string;
    descriptionId: string;
    description: string;
};

const FeatureList: FeatureItem[] = [
    {
        eyebrowId: 'homepage.features.runtime.eyebrow',
        eyebrow: 'Transport',
        titleId: 'homepage.features.runtime.title',
        title: 'One bot, two delivery modes',
        descriptionId: 'homepage.features.runtime.description',
        description:
            'The template supports polling for local development and webhook delivery for deployed environments without changing the rest of the application shape.',
    },
    {
        eyebrowId: 'homepage.features.streaming.eyebrow',
        eyebrow: 'Handlers',
        titleId: 'homepage.features.streaming.title',
        title: 'Practical bot interaction demos',
        descriptionId: 'homepage.features.streaming.description',
        description:
            'Built-in commands cover reply keyboards, inline keyboards, inline queries, callback queries, and default fallbacks that you can extend into real workflows.',
    },
    {
        eyebrowId: 'homepage.features.operations.eyebrow',
        eyebrow: 'Operations',
        titleId: 'homepage.features.operations.title',
        title: 'Service foundations included',
        descriptionId: 'homepage.features.operations.description',
        description:
            'Configuration validation, health probes, structured logging, session backends, tests, and Docker Compose files are already part of the baseline.',
    },
];

function Feature({ eyebrowId, eyebrow, titleId, title, descriptionId, description }: FeatureItem) {
    return (
        <div className={clsx('col col--4', styles.featureCol)}>
            <div className={styles.featureCard}>
                <p className={styles.featureEyebrow}>
                    <Translate id={eyebrowId}>{eyebrow}</Translate>
                </p>
                <Heading as="h3">
                    <Translate id={titleId}>{title}</Translate>
                </Heading>
                <p>
                    <Translate id={descriptionId}>{description}</Translate>
                </p>
            </div>
        </div>
    );
}

export default function HomepageFeatures(): ReactNode {
    return (
        <section className={styles.features}>
            <div className="container">
                <div className={styles.sectionHeader}>
                    <Heading as="h2">
                        <Translate id="homepage.features.section.title">What This Docs Site Covers</Translate>
                    </Heading>
                    <p>
                        <Translate id="homepage.features.section.description">
                            Architecture, bot commands, session behavior, deployment modes, and the operational details
                            needed to turn this starter into a production Telegram bot.
                        </Translate>
                    </p>
                </div>
                <div className="row">
                    {FeatureList.map((props) => (
                        <Feature key={props.title} {...props} />
                    ))}
                </div>
            </div>
        </section>
    );
}
