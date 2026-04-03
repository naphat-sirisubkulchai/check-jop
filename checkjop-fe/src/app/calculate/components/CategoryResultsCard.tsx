"use client";

import { Card } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { CategoryResult } from "@/types";
import { CategoryProgress } from "./CategoryProgress";
import { CheckCircle } from "lucide-react";

interface CategoryResultsCardProps {
  categoryResults: CategoryResult[];
}

export function CategoryResultsCard({
  categoryResults,
}: CategoryResultsCardProps) {
  const missingCategories = categoryResults.filter(
    (category) => !category.is_satisfied
  );

  return (
    <Card className="p-6 shadow-md border-gray-200">
      <h3 className="text-2xl font-bold bg-gradient-to-r from-chula-active to-pink-500 bg-clip-text text-transparent mb-4">
        Category Requirements
      </h3>
      <Tabs defaultValue="allReq">
        <TabsList aria-label="Category requirements view options" className="bg-gray-100 p-1">
          <TabsTrigger value="allReq" className="data-[state=active]:bg-white data-[state=active]:text-chula-active font-semibold">
            All Requirements
          </TabsTrigger>
          <TabsTrigger value="missing" className="data-[state=active]:bg-white data-[state=active]:text-chula-active font-semibold">
            Missing Only ({missingCategories.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="allReq" role="tabpanel" className="mt-4">
          <div className="divide-y divide-gray-200" role="list" aria-label="All category requirements">
            {categoryResults.map((category, index) => (
              <CategoryProgress key={index} category={category} />
            ))}
          </div>
        </TabsContent>

        <TabsContent value="missing" role="tabpanel" className="mt-4">
          <div className="divide-y divide-gray-200" role="list" aria-label="Missing category requirements">
            {missingCategories.length > 0 ? (
              missingCategories.map((category, index) => (
                <CategoryProgress key={index} category={category} />
              ))
            ) : (
              <div className="p-12 text-center bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl border-2 border-green-200" role="status">
                <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full mb-4">
                  <CheckCircle className="h-8 w-8 text-white" />
                </div>
                <p className="text-lg font-bold text-green-800">All category requirements are satisfied!</p>
                <p className="text-sm text-green-600 mt-2">Great job! You've completed all required categories.</p>
              </div>
            )}
          </div>
        </TabsContent>
      </Tabs>
    </Card>
  );
}
