import { useEffect, useState } from "react";
import { useConfig, useProtocols, type Config } from "@/api/config";
import { PageHeader } from "@/components/kumo/page-header/page-header";
import { ConfigForm } from "./ConfigForm";

export function ConfigPage() {
  const seedConfig = useConfig();
  const protocols = useProtocols();
  const [config, setConfig] = useState<Config>(seedConfig);

  useEffect(() => {
    setConfig(seedConfig);
  }, [seedConfig]);

  return (
    <div className="space-y-6">
      <PageHeader
        breadcrumbs={
          <nav className="flex items-center gap-2 px-3 py-2 text-sm text-kumo-fg-secondary">
            <span>Operations</span>
            <span>/</span>
            <span className="text-kumo-fg-primary">Config</span>
          </nav>
        }
        title="Configuration"
        description="Real server/shared configuration as defined by the Go config manager and `config.yml`."
      />

      <section className="rounded-2xl border border-kumo-line bg-kumo-base p-5">
        <ConfigForm
          config={config}
          protocols={protocols.protocols}
          onSaved={(nextConfig) => {
            setConfig(nextConfig);
          }}
        />
      </section>
    </div>
  );
}
