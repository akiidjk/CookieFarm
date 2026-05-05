import { SkeletonLine } from "@cloudflare/kumo/components/loader";

export function PageSkeleton() {
  return (
    <div className="space-y-4 rounded-2xl border border-kumo-line bg-kumo-base/90 p-6">
      <SkeletonLine className="h-8 w-56" />
      <SkeletonLine className="h-4 w-96 max-w-full" />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: 4 }, (_, index) => (
          <div
            key={index}
            className="space-y-3 rounded-xl border border-kumo-line bg-kumo-control p-4"
          >
            <SkeletonLine className="h-4 w-28" />
            <SkeletonLine className="h-10 w-full" />
          </div>
        ))}
      </div>
    </div>
  );
}
