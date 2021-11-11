/*
 * lakeFS API
 * lakeFS HTTP API
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package io.lakefs.clients.api.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * MergeResultSummary
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class MergeResultSummary {
  public static final String SERIALIZED_NAME_ADDED = "added";
  @SerializedName(SERIALIZED_NAME_ADDED)
  private Integer added;

  public static final String SERIALIZED_NAME_REMOVED = "removed";
  @SerializedName(SERIALIZED_NAME_REMOVED)
  private Integer removed;

  public static final String SERIALIZED_NAME_CHANGED = "changed";
  @SerializedName(SERIALIZED_NAME_CHANGED)
  private Integer changed;

  public static final String SERIALIZED_NAME_CONFLICT = "conflict";
  @SerializedName(SERIALIZED_NAME_CONFLICT)
  private Integer conflict;


  public MergeResultSummary added(Integer added) {
    
    this.added = added;
    return this;
  }

   /**
   * Get added
   * @return added
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Integer getAdded() {
    return added;
  }


  public void setAdded(Integer added) {
    this.added = added;
  }


  public MergeResultSummary removed(Integer removed) {
    
    this.removed = removed;
    return this;
  }

   /**
   * Get removed
   * @return removed
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Integer getRemoved() {
    return removed;
  }


  public void setRemoved(Integer removed) {
    this.removed = removed;
  }


  public MergeResultSummary changed(Integer changed) {
    
    this.changed = changed;
    return this;
  }

   /**
   * Get changed
   * @return changed
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Integer getChanged() {
    return changed;
  }


  public void setChanged(Integer changed) {
    this.changed = changed;
  }


  public MergeResultSummary conflict(Integer conflict) {
    
    this.conflict = conflict;
    return this;
  }

   /**
   * Get conflict
   * @return conflict
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public Integer getConflict() {
    return conflict;
  }


  public void setConflict(Integer conflict) {
    this.conflict = conflict;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    MergeResultSummary mergeResultSummary = (MergeResultSummary) o;
    return Objects.equals(this.added, mergeResultSummary.added) &&
        Objects.equals(this.removed, mergeResultSummary.removed) &&
        Objects.equals(this.changed, mergeResultSummary.changed) &&
        Objects.equals(this.conflict, mergeResultSummary.conflict);
  }

  @Override
  public int hashCode() {
    return Objects.hash(added, removed, changed, conflict);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class MergeResultSummary {\n");
    sb.append("    added: ").append(toIndentedString(added)).append("\n");
    sb.append("    removed: ").append(toIndentedString(removed)).append("\n");
    sb.append("    changed: ").append(toIndentedString(changed)).append("\n");
    sb.append("    conflict: ").append(toIndentedString(conflict)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

