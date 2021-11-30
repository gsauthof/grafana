package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

type setUserResourcePermissionTest struct {
	desc       string
	orgID      int64
	userID     int64
	actions    []string
	resource   string
	resourceID string
	seeds      []accesscontrol.SetResourcePermissionCommand
}

func TestAccessControlStore_SetUserResourcePermission(t *testing.T) {
	tests := []setUserResourcePermissionTest{
		{
			desc:       "should set resource permission for user",
			userID:     1,
			actions:    []string{"datasources:query"},
			resource:   "datasources",
			resourceID: "1",
		},
		{
			desc:       "should remove resource permission for user",
			orgID:      1,
			userID:     1,
			actions:    []string{},
			resource:   "datasources",
			resourceID: "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:query"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
		{
			desc:       "should add new resource permission for user",
			orgID:      1,
			userID:     1,
			actions:    []string{"datasources:query", "datasources:write"},
			resource:   "datasources",
			resourceID: "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:write"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			store, _ := setupTestEnv(t)

			for _, s := range test.seeds {
				_, err := store.SetUserResourcePermission(context.Background(), test.orgID, test.userID, s)
				require.NoError(t, err)
			}

			added, err := store.SetUserResourcePermission(context.Background(), test.userID, test.userID, accesscontrol.SetResourcePermissionCommand{
				Actions:    test.actions,
				Resource:   test.resource,
				ResourceID: test.resourceID,
			})

			require.NoError(t, err)
			if len(test.actions) == 0 {
				assert.Equal(t, accesscontrol.ResourcePermission{}, *added)
			} else {
				assert.Len(t, added.Actions, len(test.actions))
				assert.Equal(t, getResourceScope(test.resource, test.resourceID), added.Scope)
			}
		})
	}
}

type setTeamResourcePermissionTest struct {
	desc       string
	orgID      int64
	teamID     int64
	actions    []string
	resource   string
	resourceID string
	seeds      []accesscontrol.SetResourcePermissionCommand
}

func TestAccessControlStore_SetTeamResourcePermission(t *testing.T) {
	tests := []setTeamResourcePermissionTest{
		{
			desc:       "should add new resource permission for team",
			orgID:      1,
			teamID:     1,
			actions:    []string{"datasources:query"},
			resource:   "datasources",
			resourceID: "1",
		},
		{
			desc:       "should add new resource permission when others exist",
			orgID:      1,
			teamID:     1,
			actions:    []string{"datasources:query", "datasources:write"},
			resource:   "datasources",
			resourceID: "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:query"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
		{
			desc:       "should remove permissions for team",
			orgID:      1,
			teamID:     1,
			actions:    []string{},
			resource:   "datasources",
			resourceID: "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:query"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			store, _ := setupTestEnv(t)

			for _, s := range test.seeds {
				_, err := store.SetTeamResourcePermission(context.Background(), test.orgID, test.teamID, s)
				require.NoError(t, err)
			}

			added, err := store.SetTeamResourcePermission(context.Background(), test.orgID, test.teamID, accesscontrol.SetResourcePermissionCommand{
				Actions:    test.actions,
				Resource:   test.resource,
				ResourceID: test.resourceID,
			})

			require.NoError(t, err)
			if len(test.actions) == 0 {
				assert.Equal(t, accesscontrol.ResourcePermission{}, *added)
			} else {
				assert.Len(t, added.Actions, len(test.actions))
				assert.Equal(t, getResourceScope(test.resource, test.resourceID), added.Scope)
			}
		})
	}
}

type setBuiltinResourcePermissionTest struct {
	desc        string
	orgID       int64
	builtinRole string
	actions     []string
	resource    string
	resourceID  string
	seeds       []accesscontrol.SetResourcePermissionCommand
}

func TestAccessControlStore_SetBuiltinResourcePermission(t *testing.T) {
	tests := []setBuiltinResourcePermissionTest{
		{
			desc:        "should add new resource permission for builtin role",
			orgID:       1,
			builtinRole: "Viewer",
			actions:     []string{"datasources:query"},
			resource:    "datasources",
			resourceID:  "1",
		},
		{
			desc:        "should add new resource permission when others exist",
			orgID:       1,
			builtinRole: "Viewer",
			actions:     []string{"datasources:query", "datasources:write"},
			resource:    "datasources",
			resourceID:  "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:query"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
		{
			desc:        "should remove permissions for builtin role",
			orgID:       1,
			builtinRole: "Viewer",
			actions:     []string{},
			resource:    "datasources",
			resourceID:  "1",
			seeds: []accesscontrol.SetResourcePermissionCommand{
				{
					Actions:    []string{"datasources:query"},
					Resource:   "datasources",
					ResourceID: "1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			store, _ := setupTestEnv(t)

			for _, s := range test.seeds {
				_, err := store.SetBuiltinResourcePermission(context.Background(), test.orgID, test.builtinRole, s)
				require.NoError(t, err)
			}

			added, err := store.SetBuiltinResourcePermission(context.Background(), test.orgID, test.builtinRole, accesscontrol.SetResourcePermissionCommand{
				Actions:    test.actions,
				Resource:   test.resource,
				ResourceID: test.resourceID,
			})

			require.NoError(t, err)
			if len(test.actions) == 0 {
				assert.Equal(t, accesscontrol.ResourcePermission{}, *added)
			} else {
				assert.Len(t, added.Actions, len(test.actions))
				assert.Equal(t, getResourceScope(test.resource, test.resourceID), added.Scope)
			}
		})
	}
}

type getResourcesPermissionsTest struct {
	desc        string
	numUsers    int
	actions     []string
	resource    string
	resourceIDs []string
}

func TestAccessControlStore_GetResourcesPermissions(t *testing.T) {
	tests := []getResourcesPermissionsTest{
		{
			desc:        "should return permissions for all resource ids",
			numUsers:    3,
			actions:     []string{"datasources:query"},
			resource:    "datasources",
			resourceIDs: []string{"1", "2"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			store, sql := setupTestEnv(t)

			for _, id := range test.resourceIDs {
				seedResourcePermissions(t, store, sql, test.actions, test.resource, id, test.numUsers)
			}

			permissions, err := store.GetResourcesPermissions(context.Background(), 1, accesscontrol.GetResourcesPermissionsQuery{
				Actions:     test.actions,
				Resource:    test.resource,
				ResourceIDs: test.resourceIDs,
			})
			require.NoError(t, err)

			expectedLen := test.numUsers * len(test.resourceIDs)
			assert.Len(t, permissions, expectedLen)
		})
	}
}

func seedResourcePermissions(t *testing.T, store *AccessControlStore, sql *sqlstore.SQLStore, actions []string, resource, resourceID string, numUsers int) {
	t.Helper()
	for i := 0; i < numUsers; i++ {
		org, _ := sql.GetOrgByName("test")

		if org == nil {
			addedOrg, err := sql.CreateOrgWithMember("test", int64(i))
			require.NoError(t, err)
			org = &addedOrg
		}

		u, err := sql.CreateUser(context.Background(), models.CreateUserCommand{
			Login: fmt.Sprintf("user:%s%d", resourceID, i),
			OrgId: org.Id,
		})
		require.NoError(t, err)

		_, err = store.SetUserResourcePermission(context.Background(), 1, u.Id, accesscontrol.SetResourcePermissionCommand{
			Actions:    actions,
			Resource:   resource,
			ResourceID: resourceID,
		})
		require.NoError(t, err)
	}
}