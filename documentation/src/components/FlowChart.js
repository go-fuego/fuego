import Mermaid from "@theme/Mermaid";
import React from "react";

const code = `flowchart TB

Request -- JSON {'a':'value '} --> Deserialization

subgraph fuego

 ErrorHandler{{ErrorHandler}}

 subgraph input
	 Deserialization -- struct{A:'value '} --> InTransformation
	 InTransformation -- struct{A:'My Value'} --> Validation
 end

 Validation -- struct{A:'My Value'} --> Controller
 Controller -- struct{B:'Response!'} --> OutTransformation

 subgraph output
  OutTransformation -- struct{B:'My Response!'} --> Serialization
 end 

 Controller{{Controller}}

 click Controller "/fuego/docs/guides/controllers" "Controllers"
 click Validation "/fuego/docs/guides/validation" "Controllers"
 click InTransformation "/fuego/docs/guides/transformation" "Transformation"
 click OutTransformation "/fuego/docs/guides/transformation" "Transformation"
 click Serialization "/fuego/docs/guides/serialization" "Serialization"
 click Deserialization "/fuego/docs/guides/serialization" "Serialization"
 click ErrorHandler "/fuego/docs/guides/errors" "Error Handling"
end

ErrorHandler -- JSON{b:'Error!'} --> resp(Response)
Serialization -- JSON{b:'My Response!'} --> resp(Response)
Controller -. JSON{b:'Custom Response!'} .-> resp(Response)


`;

export function FlowChart({ selected }) {
  let style = "";
  if (selected && typeof selected === "string") {
    if (selected === "Transformation") {
      style +=
        `style InTransformation stroke:#f33,stroke-width:4px` +
        "\n" +
        `style OutTransformation stroke:#f33,stroke-width:4px`;
    } else if (selected === "Serialization") {
      style +=
        `style Serialization stroke:#f33,stroke-width:4px` +
        "\n" +
        `style Deserialization stroke:#f33,stroke-width:4px`;
    } else {
      style += `style ${selected} stroke:#f33,stroke-width:4px`;
    }
  }

  return <Mermaid value={code + style} />;
}

export default FlowChart;
