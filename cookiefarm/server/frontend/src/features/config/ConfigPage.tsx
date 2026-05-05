import { useEffect, useState } from "react";
import { Breadcrumbs } from "@cloudflare/kumo/components/breadcrumbs";
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
          <Breadcrumbs className="px-3 py-2 text-sm">
            <Breadcrumbs.Link href="/">Operations</Breadcrumbs.Link>
            <Breadcrumbs.Separator />
            <Breadcrumbs.Current>Config</Breadcrumbs.Current>
          </Breadcrumbs>
        }
        title="Configuration"
        description="Server/Shared configuration as defined by the Go config manager and `config.yml`."
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
