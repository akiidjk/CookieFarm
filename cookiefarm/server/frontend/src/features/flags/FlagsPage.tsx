import { useMemo, useState } from "react";
import { Banner } from "@cloudflare/kumo/components/banner";
import { Breadcrumbs } from "@cloudflare/kumo";
import { Button } from "@cloudflare/kumo/components/button";
import { Select } from "@cloudflare/kumo/components/select";
import { CaretLeft, CaretRight, WarningCircle } from "@phosphor-icons/react";
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
  cursor: string,
  pageSize: number,
  filters: FlagFilterState,
): FlagsQuery {
  return {
    limit: pageSize,
    ...(cursor ? { cursor } : {}),
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
  const [cursor, setCursor] = useState<string>("");
  const [previousCursors, setPreviousCursors] = useState<string[]>([]);
  const [pageSize, setPageSize] = useState(defaultPageSize);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const debouncedSearch = useDebounce(filters.search, 300);
  const flagsRequest = useMemo(
    () =>
      buildFlagsRequest(cursor, pageSize, {
        ...filters,
        search: debouncedSearch,
      }),
    [cursor, debouncedSearch, filters, pageSize],
  );
  const flagsQuery = useFlags(flagsRequest, {
    refreshInterval: flagsRefreshInterval,
  });
  const rows = flagsQuery.data!.flags;
  const totalCount = flagsQuery.data!.n_flags;
  const nextCursor = flagsQuery.data!.next;

  const currentPage = previousCursors.length + 1;
  const totalPages = Math.max(1, Math.ceil(totalCount / pageSize));

  function resetCursor() {
    setCursor("");
    setPreviousCursors([]);
  }

  function goToNextPage() {
    if (!nextCursor) return;
    setPreviousCursors((prev) => [...prev, cursor]);
    setCursor(nextCursor);
  }

  function goToPrevPage() {
    if (previousCursors.length === 0) return;
    const updated = [...previousCursors];
    const prevCursor = updated.pop()!;
    setPreviousCursors(updated);
    setCursor(prevCursor);
  }

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
          resetCursor();
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
              resetCursor();
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

      <section className="flex items-center justify-between rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <span className="text-sm text-kumo-fg-secondary">
          Page {currentPage} of {totalPages} &mdash; {totalCount} total flags
        </span>
        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            onClick={goToPrevPage}
            disabled={previousCursors.length === 0}
          >
            <CaretLeft weight="bold" />
            Previous
          </Button>
          <Button
            variant="secondary"
            onClick={goToNextPage}
            disabled={!nextCursor}
          >
            Next
            <CaretRight weight="bold" />
          </Button>
        </div>
      </section>
    </div>
  );
}
