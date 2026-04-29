import { useMemo, useState } from "react";
import { Banner } from "@cloudflare/kumo/components/banner";
import { Breadcrumbs } from "@cloudflare/kumo";
import { Button } from "@cloudflare/kumo/components/button";
import { Pagination } from "@cloudflare/kumo/components/pagination";
import { Select } from "@cloudflare/kumo/components/select";
import { WarningCircle } from "@phosphor-icons/react";
import { useConfig } from "@/api/config";
import { useFlags, type FlagStatus, type FlagsQuery } from "@/api/flags";
import { PageHeader } from "@/components/kumo/page-header/page-header";
import { useDebounce } from "@/hooks/useDebounce";
import { FlagFilters, type FlagFilterState } from "./FlagFilters";
import { FlagTable } from "./FlagTable";

const defaultPageSize = 40;
const flagsRefreshInterval = 10_000;

const defaultFilters: FlagFilterState = {
  status: "all",
  service: "",
  team: "",
  search: "",
  searchField: "flag_code",
};

function buildFlagsRequest(
  page: number,
  pageSize: number,
  filters: FlagFilterState,
): FlagsQuery {
  return {
    limit: pageSize,
    offset: (page - 1) * pageSize,
    ...(filters.status !== "all" ? { status: Number(filters.status) as FlagStatus } : {}),
    ...(filters.service ? { service: filters.service } : {}),
    ...(filters.team.trim() ? { team: filters.team.trim() } : {}),
    ...(filters.search.trim() ? { search: filters.search.trim() } : {}),
    ...(filters.searchField ? { searchField: filters.searchField } : {}),
  };
}

export function FlagsPage() {
  const config = useConfig();
  const [filters, setFilters] = useState<FlagFilterState>(defaultFilters);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(defaultPageSize);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const debouncedSearch = useDebounce(filters.search, 300);
  const flagsRequest = useMemo(
    () =>
      buildFlagsRequest(page, pageSize, {
        ...filters,
        search: debouncedSearch,
      }),
    [debouncedSearch, filters, page, pageSize],
  );
  const flagsQuery = useFlags(flagsRequest, {
    refreshInterval: flagsRefreshInterval,
  });
  const rows = flagsQuery.data!.flags;
  const totalCount = flagsQuery.data!.n_flags;

  async function refreshFlags(): Promise<void> {
    await flagsQuery.mutate();
    setErrorMessage(null);
  }
  const visibleErrorMessage =
    errorMessage ?? (flagsQuery.error instanceof Error ? flagsQuery.error.message : null);

  return (
    <div className="space-y-6">
      <PageHeader
        breadcrumbs={
          <Breadcrumbs className="px-3 py-2">
            <Breadcrumbs.Link href="/">Operations</Breadcrumbs.Link>
            <Breadcrumbs.Separator />
            <Breadcrumbs.Current>Flags</Breadcrumbs.Current>
          </Breadcrumbs>
        }
        title="Flags"
        description="Handle flag management and filtering"
      >
        <Button
          variant="secondary"
          onClick={() => {
            void refreshFlags().catch((error: unknown) => {
              setErrorMessage(error instanceof Error ? error.message : "Flag refresh failed");
            });
          }}
        >
          Reload
        </Button>
      </PageHeader>

      {visibleErrorMessage ? (
        <Banner
          variant="error"
          icon={<WarningCircle weight="fill" />}
          title="Unable to refresh flags"
          description={visibleErrorMessage}
        />
      ) : null}

      <FlagFilters
        filters={filters}
        services={config.shared.services}
        onChange={(nextFilters) => {
          setPage(1);
          setFilters(nextFilters);
        }}
      />

      <section className="flex flex-wrap items-center gap-3 rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <div className="flex items-center gap-3">
          <span className="text-sm text-kumo-fg-secondary">Rows per page</span>
          <Select
            aria-label="Rows per page"
            value={String(pageSize)}
            onValueChange={(value) => {
              setPage(1);
              setPageSize(Number(value));
            }}
            items={{
              20: "20",
              40: "40",
              80: "80",
              100: "100",
            }}
          />
        </div>
      </section>

      <FlagTable rows={rows} />

      <section className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <Pagination
          page={page}
          setPage={setPage}
          perPage={pageSize}
          totalCount={totalCount}
        >
          <Pagination.Info />
          <Pagination.Separator />
          <Pagination.Controls />
        </Pagination>
      </section>
    </div>
  );
}
