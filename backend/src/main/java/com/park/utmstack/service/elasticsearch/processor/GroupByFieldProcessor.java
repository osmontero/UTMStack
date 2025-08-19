package com.park.utmstack.service.elasticsearch.processor;

import java.util.*;
import java.util.stream.Collectors;

public class GroupByFieldProcessor implements SearchResultProcessor {

    private final String fieldName;

    public GroupByFieldProcessor(String fieldName) {
        this.fieldName = fieldName;
    }

    @Override
    public List<Map<String, Object>> process(List<Map<String, Object>> rawResults) {

        Map<Object, Map<String, Object>> byId = new LinkedHashMap<>();
        for (Map<String, Object> item : rawResults) {
            Object id = item.get("id");
            if (id != null) {
                byId.put(id, item);
            }
        }

        Map<Object, List<Map<String, Object>>> childrenByParent = new LinkedHashMap<>();
        for (Map<String, Object> item : rawResults) {
            Object parentId = item.get("parentId");
            if (parentId != null && byId.containsKey(parentId)) {
                childrenByParent.computeIfAbsent(parentId, k -> new ArrayList<>()).add(item);
            }
        }

        List<Map<String, Object>> result = new ArrayList<>();
        for (Map<String, Object> item : rawResults) {
            Object id = item.get("id");
            if (id != null && childrenByParent.containsKey(id)) {
                Map<String, Object> enriched = new LinkedHashMap<>(item);
                List<Map<String, Object>> children = childrenByParent.get(id);
                enriched.put("children", children);
                enriched.put("groupSize", children.size());
                enriched.put("hasChildren", true);
                result.add(enriched);
            } else {
                Map<String, Object> untouched = new LinkedHashMap<>(item);
                untouched.put("children", Collections.emptyList());
                untouched.put("hasChildren", false);
                result.add(untouched);
            }
        }

        return result;
    }
}

