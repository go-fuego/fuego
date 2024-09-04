import clsx from "clsx";
import styles from "./styles.module.css";

export default function HomepageFeatures(): JSX.Element {
  return (
    <div className={styles.features} style={{}}>
      <div id="simple" className={clsx(styles.block, styles.block_odd)}>
        <h2 className="">Serialization</h2>
        <p>Simple serialization of data structures to JSON and back.</p>
      </div>
      <div id="openapi" className={clsx(styles.block)}>
        <h2 className="">OpenAPI generation</h2>
        <p>Generate OpenAPI 3.0 specifications from your data structures.</p>
      </div>
      <div id="validation" className={clsx(styles.block, styles.block_odd)}>
        <h2 className="">Validation</h2>
        <p>Easily validate data structures against your own rules.</p>
      </div>
    </div>
  );
}
