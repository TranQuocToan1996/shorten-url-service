"use client";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";

export default function SwaggerPage() {
  return (
    <div style={{ height: '100vh', background: '#fff', padding: '1em' }}>
      <SwaggerUI url="/swagger.yaml" />
    </div>
  );
}
