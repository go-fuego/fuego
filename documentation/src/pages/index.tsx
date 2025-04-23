import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import Heading from "@theme/Heading";

import styles from "./index.module.css";

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx(styles.heroBanner)}>
      <div className="container">
        <img src="/img/logo.svg" alt="" width={200} height={200} />

        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
      </div>
    </header>
  );
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <HomepageHeader />
      <main className={styles.main}>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs/">
            Tutorial - 5 min ⏱️
          </Link>
        </div>
        <iframe
          className={styles.video}
          width="640"
          height="360"
          src="https://www.youtube.com/embed/DqU7f1IKU1g?si=F6KwY5Zmh8FxCDXI&amp;start=804"
          title="Introducing Fuego !"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          referrerpolicy="strict-origin-when-cross-origin"
          allowfullscreen
        ></iframe>
      </main>
    </Layout>
  );
}
