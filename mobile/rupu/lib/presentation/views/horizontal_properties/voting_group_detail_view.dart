import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rupu/domain/entities/voting.dart';
import 'package:rupu/domain/entities/voting_results_result.dart';
import 'package:rupu/presentation/views/horizontal_properties/controllers/voting_group_detail_controller.dart';
import 'package:rupu/presentation/views/horizontal_properties/horizontal_property_detail_view.dart';

class VotingGroupDetailView extends StatelessWidget {
  final int propertyId;
  final int votingGroupId;

  const VotingGroupDetailView({
    super.key,
    required this.propertyId,
    required this.votingGroupId,
  });

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<VotingGroupDetailController>();
    return Scaffold(
      appBar: AppBar(
        title: Obx(() => Text(controller.group.value?.name ?? 'Detalle de grupo')),
      ),
      body: Obx(
        () {
          if (controller.isLoading.value) {
            return const Center(child: CircularProgressIndicator());
          }
          if (controller.errorMessage.value != null) {
            return Center(child: Text(controller.errorMessage.value!));
          }
          return ListView.builder(
            itemCount: controller.votings.length,
            itemBuilder: (context, index) {
              final voting = controller.votings[index];
              return _VotingCard(voting: voting);
            },
          );
        },
      ),
    );
  }
}

class _VotingCard extends StatelessWidget {
  final Voting voting;

  const _VotingCard({required this.voting});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;

    return Card(
      margin: const EdgeInsets.all(8.0),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(voting.title, style: tt.titleMedium),
            const SizedBox(height: 8.0),
            Text(voting.description, style: tt.bodyMedium),
            const SizedBox(height: 8.0),
            Row(
              children: [
                Switch(
                  value: voting.isActive,
                  onChanged: (value) {
                    if (value) {
                      Get.find<VotingGroupDetailController>()
                          .activateVoting(voting.id);
                    } else {
                      Get.find<VotingGroupDetailController>()
                          .deactivateVoting(voting.id);
                    }
                  },
                ),
                Text(voting.isActive ? 'Activa' : 'Inactiva'),
                const Spacer(),
                _CardActions(
                  onLive: () {
                    final state = GoRouter.of(context).routerDelegate.currentConfiguration;
                    final segments = state.uri.pathSegments;
                    final page = segments.length > 1 ? segments[1] : '0';
                    final propertyId =
                        Get.find<VotingGroupDetailController>().propertyId;
                    final path =
                        '/home/$page/horizontal-properties/$propertyId/voting/${voting.votingGroupId}/votings/${voting.id}/live';
                    GoRouter.of(context).push(path);
                  },
                  onEdit: () {},
                  onDelete: () {},
                  onMore: () {},
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
